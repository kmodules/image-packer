# Test scripts

```bash
docker network create registry
docker run -d --restart=unless-stopped --name registry --net registry \
  -e "REGISTRY_STORAGE_FILESYSTEM_ROOTDIRECTORY=/var/lib/registry" \
  -e "REGISTRY_STORAGE_DELETE_ENABLED=true" \
  -e "REGISTRY_VALIDATION_DISABLED=true" \
  -v "registry-data:/var/lib/registry" \
  -p "127.0.0.1:5000:5000" \
  registry:2
```

```bash
image-packer generate-scripts --output-dir=catalog \
    --src=catalog/imagelist.yaml
```
