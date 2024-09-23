/*
Copyright AppsCode Inc. and Contributors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmds

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"kmodules.xyz/client-go/tools/parser"

	"github.com/spf13/cobra"
	"gomodules.xyz/sets"
	apps "k8s.io/api/apps/v1"
	batch "k8s.io/api/batch/v1"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	api "kubedb.dev/apimachinery/apis/catalog/v1alpha1"
)

func NewCmdListImages() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "list",
		Short:                 "List all Docker images in a dir/file or stdin",
		DisableFlagsInUseLine: true,
		DisableAutoGenTag:     true,
		RunE: func(cmd *cobra.Command, args []string) error {
			images, err := listImages(args)
			if err != nil {
				return err
			}
			fmt.Println(strings.Join(images, "\n"))
			return nil
		},
	}

	return cmd
}

func listImages(args []string) ([]string, error) {
	dir, manifest, err := readManifest(args)
	if err != nil {
		return nil, err
	}
	return ListImages(dir, manifest)
}

func readManifest(args []string) (string, bool, error) {
	if len(args) == 0 {
		return "", false, errors.New("missing input")
	} else if len(args) > 1 {
		return "", false, errors.New("too many inputs")
	}
	dir := args[0]
	if dir == "-" {
		// ref: https://gist.github.com/AlexMocioi/10008287
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			return "", false, err
		}
		return string(data), true, nil
	}
	return dir, false, nil
}

func ListImages(dir string, manifest bool) ([]string, error) {
	imgList := sets.NewString()
	if manifest {
		err := parser.ProcessResources([]byte(dir), processYAML(imgList))
		if err != nil {
			return nil, err
		}
	} else {
		err := parser.ProcessPath(dir, processYAML(imgList))
		if err != nil {
			return nil, err
		}
	}
	return imgList.List(), nil
}

func processYAML(imgList sets.String) func(ri parser.ResourceInfo) error {
	return func(ri parser.ResourceInfo) error {
		switch ri.Object.GetObjectKind().GroupVersionKind() {
		case core.SchemeGroupVersion.WithKind("ReplicationController"):
			var obj core.ReplicationController
			if err := runtime.DefaultUnstructuredConverter.FromUnstructured(ri.Object.UnstructuredContent(), &obj); err != nil {
				return err
			}
			collectFromContainers(obj.Spec.Template.Spec.Containers, imgList)
			collectFromContainers(obj.Spec.Template.Spec.InitContainers, imgList)
		case apps.SchemeGroupVersion.WithKind("Deployment"):
			var obj apps.Deployment
			if err := runtime.DefaultUnstructuredConverter.FromUnstructured(ri.Object.UnstructuredContent(), &obj); err != nil {
				return err
			}
			collectFromContainers(obj.Spec.Template.Spec.Containers, imgList)
			collectFromContainers(obj.Spec.Template.Spec.InitContainers, imgList)
		case apps.SchemeGroupVersion.WithKind("StatefulSet"):
			var obj apps.StatefulSet
			if err := runtime.DefaultUnstructuredConverter.FromUnstructured(ri.Object.UnstructuredContent(), &obj); err != nil {
				return err
			}
			collectFromContainers(obj.Spec.Template.Spec.Containers, imgList)
			collectFromContainers(obj.Spec.Template.Spec.InitContainers, imgList)
		case apps.SchemeGroupVersion.WithKind("DaemonSet"):
			var obj apps.DaemonSet
			if err := runtime.DefaultUnstructuredConverter.FromUnstructured(ri.Object.UnstructuredContent(), &obj); err != nil {
				return err
			}
			collectFromContainers(obj.Spec.Template.Spec.Containers, imgList)
			collectFromContainers(obj.Spec.Template.Spec.InitContainers, imgList)
		case batch.SchemeGroupVersion.WithKind("Job"):
			var obj batch.Job
			if err := runtime.DefaultUnstructuredConverter.FromUnstructured(ri.Object.UnstructuredContent(), &obj); err != nil {
				return err
			}
			collectFromContainers(obj.Spec.Template.Spec.Containers, imgList)
			collectFromContainers(obj.Spec.Template.Spec.InitContainers, imgList)
		case batch.SchemeGroupVersion.WithKind("CronJob"):
			var obj batch.CronJob
			if err := runtime.DefaultUnstructuredConverter.FromUnstructured(ri.Object.UnstructuredContent(), &obj); err != nil {
				return err
			}
			collectFromContainers(obj.Spec.JobTemplate.Spec.Template.Spec.Containers, imgList)
			collectFromContainers(obj.Spec.JobTemplate.Spec.Template.Spec.InitContainers, imgList)
		case api.SchemeGroupVersion.WithKind(api.ResourceKindElasticsearchVersion):
			var v api.ElasticsearchVersion
			err := runtime.DefaultUnstructuredConverter.FromUnstructured(ri.Object.UnstructuredContent(), &v)
			if err != nil {
				return err
			}

			collect(v.Spec.DB.Image, imgList)
			collect(v.Spec.InitContainer.Image, imgList)
			collect(v.Spec.Exporter.Image, imgList)
			collect(v.Spec.Dashboard.Image, imgList)
			collect(v.Spec.DashboardInitContainer.YQImage, imgList)
		case api.SchemeGroupVersion.WithKind(api.ResourceKindMemcachedVersion):
			var v api.MemcachedVersion
			err := runtime.DefaultUnstructuredConverter.FromUnstructured(ri.Object.UnstructuredContent(), &v)
			if err != nil {
				return err
			}

			collect(v.Spec.DB.Image, imgList)
			collect(v.Spec.Exporter.Image, imgList)
		case api.SchemeGroupVersion.WithKind(api.ResourceKindMariaDBVersion):
			var v api.MariaDBVersion
			err := runtime.DefaultUnstructuredConverter.FromUnstructured(ri.Object.UnstructuredContent(), &v)
			if err != nil {
				return err
			}

			collect(v.Spec.DB.Image, imgList)
			collect(v.Spec.Exporter.Image, imgList)
			collect(v.Spec.InitContainer.Image, imgList)
			collect(v.Spec.Coordinator.Image, imgList)
		case api.SchemeGroupVersion.WithKind(api.ResourceKindMongoDBVersion):
			var v api.MongoDBVersion
			err := runtime.DefaultUnstructuredConverter.FromUnstructured(ri.Object.UnstructuredContent(), &v)
			if err != nil {
				return err
			}

			collect(v.Spec.DB.Image, imgList)
			collect(v.Spec.Exporter.Image, imgList)
			collect(v.Spec.InitContainer.Image, imgList)
			collect(v.Spec.ReplicationModeDetector.Image, imgList)
		case api.SchemeGroupVersion.WithKind(api.ResourceKindMySQLVersion):
			var v api.MySQLVersion
			err := runtime.DefaultUnstructuredConverter.FromUnstructured(ri.Object.UnstructuredContent(), &v)
			if err != nil {
				return err
			}

			collect(v.Spec.DB.Image, imgList)
			collect(v.Spec.Exporter.Image, imgList)
			collect(v.Spec.InitContainer.Image, imgList)
			collect(v.Spec.ReplicationModeDetector.Image, imgList)
			collect(v.Spec.Coordinator.Image, imgList)
			collect(v.Spec.Router.Image, imgList)
			collect(v.Spec.RouterInitContainer.Image, imgList)
		case api.SchemeGroupVersion.WithKind(api.ResourceKindPerconaXtraDBVersion):
			var v api.PerconaXtraDBVersion
			err := runtime.DefaultUnstructuredConverter.FromUnstructured(ri.Object.UnstructuredContent(), &v)
			if err != nil {
				return err
			}

			collect(v.Spec.DB.Image, imgList)
			collect(v.Spec.Exporter.Image, imgList)
			collect(v.Spec.InitContainer.Image, imgList)
		case api.SchemeGroupVersion.WithKind(api.ResourceKindPgBouncerVersion):
			var v api.PgBouncerVersion
			err := runtime.DefaultUnstructuredConverter.FromUnstructured(ri.Object.UnstructuredContent(), &v)
			if err != nil {
				return err
			}

			collect(v.Spec.PgBouncer.Image, imgList)
			collect(v.Spec.Exporter.Image, imgList)
		case api.SchemeGroupVersion.WithKind(api.ResourceKindProxySQLVersion):
			var v api.ProxySQLVersion
			err := runtime.DefaultUnstructuredConverter.FromUnstructured(ri.Object.UnstructuredContent(), &v)
			if err != nil {
				return err
			}

			collect(v.Spec.Proxysql.Image, imgList)
			collect(v.Spec.Exporter.Image, imgList)
		case api.SchemeGroupVersion.WithKind(api.ResourceKindRedisVersion):
			var v api.RedisVersion
			err := runtime.DefaultUnstructuredConverter.FromUnstructured(ri.Object.UnstructuredContent(), &v)
			if err != nil {
				return err
			}

			collect(v.Spec.DB.Image, imgList)
			collect(v.Spec.Exporter.Image, imgList)
			collect(v.Spec.Coordinator.Image, imgList)
			collect(v.Spec.InitContainer.Image, imgList)
		case api.SchemeGroupVersion.WithKind(api.ResourceKindPostgresVersion):
			var v api.PostgresVersion
			err := runtime.DefaultUnstructuredConverter.FromUnstructured(ri.Object.UnstructuredContent(), &v)
			if err != nil {
				return err
			}

			collect(v.Spec.DB.Image, imgList)
			collect(v.Spec.Coordinator.Image, imgList)
			collect(v.Spec.Exporter.Image, imgList)
			collect(v.Spec.InitContainer.Image, imgList)
		}
		return nil
	}
}

func collect(ref string, dm sets.String) {
	if ref == "" {
		return
	}
	dm.Insert(ref)
}

func collectFromContainers(containers []core.Container, dm sets.String) {
	for _, c := range containers {
		dm.Insert(c.Image)
	}
}
