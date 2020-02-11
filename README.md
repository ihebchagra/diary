# diary

## Dependencies

* vim or neovim for editing text
* fzf for searching through entries

## Usage

You can set where you want to save your diary entries with the environmental variable $DIARY_DIR.
By default this will be set to ~/dox/diary for my personal preferences

* you can start writing in today's entry with
```
$ diary
```

* you can search through older entries with
```
$ diary search
```
## Building:

```
$ go build main.go
```

