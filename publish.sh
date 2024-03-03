Version="v0.0.13"
# projDir="./cd-user/user.go"

# cd $projDir
go mod tidy
git submodule update --remote
git add cd.go go.mod go.sum sys/base/b.go sys/base/cd-error.go sys/base/go sys/base/go.mod  sys/base/go.sum
git add -A
git commit -a -m "set version $Version"
git tag $Version
git push origin $Version

# cd-cli mod publish 

