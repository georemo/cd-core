# get current repository latest version
echo "current repository latest version:\n"
git ls-remote --tags https://github.com/tcp-x/cd-core.git
# set latest version
Version="v0.0.58"
# projDir="./cd-user/user.go"

# cd $projDir
go mod tidy
git submodule update --remote
git add sys/base/
git commit -am 'Add package github.com/tcp-x/cd-core/sys/base'
git tag $Version-base
git push
# git add cd.go go.mod go.sum sys/base/b.go sys/base/IBase.go sys/base/cd-error.go sys/base/go.mod  sys/base/go.sum sys/user/user.go sys/user/session.go sys/user/go.mod  sys/user/go.sum
git add sys/user/
git commit -am 'Add package github.com/tcp-x/cd-core/sys/user'
git tag $Version-user
git push

git add .
git commit -a -m "set version $Version"
git tag $Version
git push origin $Version


# cd-cli mod publish 

