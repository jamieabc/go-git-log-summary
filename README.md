This repo is to summarize git commit behavior.

For some reason, git might pack data to save disk space, and it will cause parsing additional overhead. My goal is to find better statics of specific directory but not learning internal of git, so I find another way to solve the problem of git packing data.

At root of repository, type following commands to unpack git objects:

```shell script
git gc
mv .git/objects/pack .
git unpack-objects < pack/*.pack
rm -rf pack
```

reference about git internals
1. [git object model](http://shafiul.github.io/gitbook/1_the_git_object_model.html)
1. [git internals](https://git-scm.com/book/id/v2/Git-Internals-Git-Objects)
