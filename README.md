# pgtable

A modern pager of displaying results of `psql`.

## Installation

Install the application with go:

```sh
go install github.com/untanky/pgtable/cmd/pgtable@0.1.0
```

To use it as a pager, you need point the `PAGER` variable to the binary. If `$GOBIN` is in your path, adding this to your `.psqlrc` suffices:

```psql
…
\setenv PAGER pgtable
…
```

You can also set the environment variable before invoking `psql`:

```sh
PAGER=pgtable psql
```

> The feature set is currently limited, therefore some adjustments to the `.psqlrc` are necessary to get this working properly.
>
> ```psql
> \timing off -- this is the default setting
> 
> \pset null '[null]' -- enable pgtable to detect null
> \pset border 1 -- this is the default setting
> \pset linestyle ascii -- this is the default setting
> ```

## Keybindings:

Keys | Action
--- | ---
`<ctrl>+c/q` | Quit the pager
`h`/`<left>` | Move left
`l`/`<right>` | Move right
`j`/`<down>` | Move down
`k`/`<up>` | Move up
`<ctrl>+d` | Move down half a screen height
`<ctrl>+u` | Move up half a screen height
`yy` | Move the current cell value to the pasteboard

