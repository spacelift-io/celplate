# celplate

[![Test](https://github.com/spacelift-io/celplate/actions/workflows/test.yml/badge.svg)](https://github.com/spacelift-io/celplate/actions/workflows/test.yml)

Celplate provides an elegant way to template files using the [Common Expression Language (CEL)](https://github.com/google/cel-spec).
The package comes with all the batteries included: a generic [scanner](scanner.go) and a single (CEL) [evaluator](evaluator/cel.go) for it.

By default the scanner assumes that any block of text that start with `${{` is a special input block 
which has to be evaluated and must end with `}}`, example:

``` yaml
${{ inputs.serial }}
```

Too see the library in action, checkout the [end to end test](e2e) for it.

**Note that the current implementation does not support escaping the special input block.**

## Extensions

We've added [ext/Strings](https://pkg.go.dev/github.com/google/cel-go/ext#Strings) extensions to the CEL evaluator. This includes a bunch of useful methods, such as `charAt`, `indexOf`, `join`, `split`, `replace`, `trim` etc.

## Releasing

To create a new release just create a new tag with a `v` prefix and push it to main. For more details checkout the go [docs on publishing modules](https://go.dev/blog/publishing-go-modules).
