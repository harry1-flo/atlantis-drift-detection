# Drift-Detection

## Parameters

| Flag                    | Short | Description                     | Example                        |
| ----------------------- | ----- | ------------------------------- | ------------------------------ |
| `--github-token`        | `-g`  | The Github token                | `ghp_xxx`                      |
| `--github-repo-ref`     | `-f`  | The Github repository reference | `main`, `master`               |
| `--atlantis-url`        | `-u`  | The Atlantis URL                | `https://atlantis.example.com` |
| `--atlantis-token`      | `-t`  | The Atlantis token              | `your-api-secret`              |
| `--atlantis-repository` | `-r`  | Atlantis Repository             | `owner/repo-name`              |
| `--atlantis-config`     | `-c`  | Atlantis Config File            | `atlantis.yaml`                |
| `--slack-bot-token`     | `-s`  | Slack Bot Token                 | `xoxb-xxx`                     |
| `--slack-channel`       | `-l`  | Slack Channel                   | `C024BE91L`                    |

## Execute

```sh
    make build
    make run
```

## Release Update

```sh
    brew install gh

    gh auth login

    gh relase create [version] --title [title] --notes [notes] path
```

./tt notification \
--at-branch-ref tt \
--at-branch-name tt \
--at-repo-name tt \
--at-commit-hash tt \
--at-pr-num tt \
--at-pr-url tt \
--at-pr-author tt \
--at-gh-token tt \
--at-command validate \
--at-owner tt \
--at-repo-rel-dir tt \
--at-slack-bottoken tt \
--at-slack-channel tt \
--at-outputs asdf
