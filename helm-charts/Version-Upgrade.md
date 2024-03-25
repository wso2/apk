# Version Upgrade APK

helm template test . -f version-upgrade-values.yaml && helm show crds . > t.yaml