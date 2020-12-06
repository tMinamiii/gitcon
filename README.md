# gitcon

.git/configを操作するツール

## 使い方

### レポジトリURLのssh化

```sh
gitcon ssh
```

```config
[remote "origin"]
    url = git@github.com:tMinamiii/gitcon.git
    fetch = +refs/heads/*:refs/remotes/origin/*
```

### レポジトリURLのhttps化

```sh
gitcon https
```

```config
[remote "origin"]
    url = https://github.com/tMinamiii/gitcon
    fetch = +refs/heads/*:refs/remotes/origin/*
```

