# ✨ Datasets

A simple tool for reproducable dataset gathering for machine learning tasks.

## Installation

MacOS/Linux
```bash
curl https://raw.githubusercontent.com/ex3ndr/datasets/main/install.sh | sh
```

## How to use

Entry point is the `datasets.yaml` file in your project, that looks like this:

```yaml
datasets:
  - cifar-100 # This downloads from central repository
  - name: some_private_dataset
    source: https://not-so-real-url.org
```

Then you can execute sync of datasets:
```bash
datasets sync
```

That's all!

## License

MIT
