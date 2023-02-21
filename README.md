# go-matcher - graphite glob/seriesByTag expressions batch match engine for Go
go-matcher is a graphite glob/seriesByTag expressions batch match engine for Go.
It doesn't have constant time guarantees like the built-in `regexp` package, but it allows backtracking and is compatible with `regexp` package.

## Basis of the engine
Contains 2 parts:
* gglob - graphite glob expressions match engine
* gtags - graphite tags (with seriesByTag) expressions match engine

## Installing
This is a go-gettable library, so install is easy:

    go get github.com/msaf1980/go-matcher/gglob
    go get github.com/msaf1980/go-matcher/gtags

## Usage

### gglob

```go
  w:= gglob.NewGlobMatcher()
  err = w.Adds([]string{"a.b", "d.*"})
  if err != nil {
    ...
  }
  
  // get matched globs (also contain normalized versions of globs)
  matchedGlobs := w.Match(path)
  
  // use preallocated slice (for better perfomance
  // var matchedGlobs []string
  matchedGlobs := matchedGlobs[:0]
  w.MatchB(path, &matchedGlobs)
```

With splitted path parts
```go
  w:= gglob.NewGlobMatcher()
  err = w.Adds([]string{"a.b", "d.*"})
  if err != nil {
    ...
  }
  
  parts := wildcards.PathSplit(path)
  
  // get matched globs (also contain normalized versions of globs)
  matchedGlobs := w.MatchByParts(parts)
  
  // use preallocated slice (for better perfomance
  // var matchedGlobs []string
  matchedGlobs := matchedGlobs[:0]
  w.MatchByPartsB(parts, &matchedGlobs)
```


Get macthed globs index
```go

  var buf strings.Builder
  buf.Grow(128)

  globs := []string{"a.b", "d.*"}
  w:= gglob.NewGlobMatcher()
  for i, glob := range globs {
    _, err = w.AddIndexed(glob, i, &buf)
    if err != nil {
      ...
    }
  }
    
  // get matched globs indexes
  matchedIndex := w.MatchIndexed(path)
  
  // use preallocated slice (for better perfomance)
  // var matchedIndex []int
  matchedIndex := matchedIndex[:0]
  w.MatchB(path, &matchedIndex)  
  
  // get first mached (with lowest index number)
  first := -1
  w.MatchFirst(path, &first)
```


Get macthed globs index with splitted path parts
```go

  var buf strings.Builder
  buf.Grow(128)

  globs := []string{"a.b", "d.*"}
  w:= gglob.NewGlobMatcher()
  for i, glob := range globs {
    _, err = w.AddIndexed(glob, i, &buf)
    if err != nil {
      ...
    }
  }
  
  parts := wildcards.PathSplit(path)
  
  // get matched globs indexes
  matchedIndex := w.MatchIndexedByParts(parts)
  
  // use preallocated slice (for better perfomance)
  // var matchedIndex []int
  matchedIndex := matchedIndex[:0]
  w.MatchByPartsB(parts, &matchedIndex)  
  
  // get first mached (with lowest index number)
  first := -1
  w.MatchFirstByParts(parts, &first)
```

### gtags


```go
  w:= gtags.NewTagsMatcher()
  err = w.Adds([]string{`seriesByTag('__name__=a.b','b=d.*', '__name__=a.b','b=e')`})
  if err != nil {
    ...
  }
  
  tags, err := PathTags(path)
  if err != nil {
    ..
  }
  
  // get matched globs (also contain normalized versions of globs)
  matchedQueries := w.MatchByTags(tags)
  
  // use preallocated slice (for better perfomance
  // var matchedQueries []string
  matchedQueries := matchedQueries[:0]
  w.MatchByTagsB(tags, &matchedQueries)
```

Get macthed globs index
```go

  var buf strings.Builder
  buf.Grow(128)

  queries := []string{`seriesByTag('__name__=a.b','b=d.*', '__name__=a.b','b=e')`}
  w:= gtags.NewTagsMatcher()
  for i, query := range queries {
    _, err = w.AddIndexed(glob, i, &buf)
    if err != nil {
      ...
    }
  }
  
  // get matched globs indexes
  matchedIndex := w.MatchIndexedByTags(path)
  
  // use preallocated slice (for better perfomance)
  // var matchedIndex []int
  matchedIndex := matchedIndex[:0]
  w.MatchByTagsB(path, &matchedIndex)  
  
  // get first mached (with lowest index number)
  first := -1
  w.MatchFirstByTags(path, &first)
```
