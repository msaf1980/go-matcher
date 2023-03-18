package gtags

import (
	"math/rand"
	"regexp"
	"strings"

	"github.com/google/go-cmp/cmp"
)

var (
	cmpTransform = cmp.Transformer("Re", func(in *regexp.Regexp) string {
		if in == nil {
			return "<nil>"
		}
		return in.String()
	})
)

func storedTagsList(paths []string) (list [][]Tag) {
	var err error
	list = make([][]Tag, len(paths))
	for i, path := range paths {
		if list[i], err = PathTags(path); err != nil {
			panic(err)
		}
	}
	return
}

func tagsList(paths []string) (list [][]Tag) {
	var err error
	list = make([][]Tag, len(paths))
	for i, path := range paths {
		if list[i], err = GraphitePathTags(path); err != nil {
			panic(err)
		}
	}
	return
}

func tagMapList(paths []string) (list []map[string]string) {
	var err error
	list = make([]map[string]string, len(paths))
	for i, path := range paths {
		if list[i], err = GraphitePathTagsMap(path); err != nil {
			panic(err)
		}
	}
	return
}

func taggedTermListList(queries []string) (list []TaggedTermList) {
	var err error
	list = make([]TaggedTermList, len(queries))
	for i, query := range queries {
		if list[i], err = ParseSeriesByTag(query); err != nil {
			panic(err)
		}
	}

	return
}

func generateTaggedMetrics(termsList []TaggedTermList, count int) []string {
	result := make([]string, 0, 3*count)

	for i := 0; i < count; i++ {
		terms := termsList[rand.Intn(len(termsList))]

		matchedMetric := generateMatchedMetric(terms)
		partiallyMatchedMetric := generatePartiallyMatchedMetric(terms)
		notMatchedMetric := generateNotMatchedMetric(terms)

		result = append(result, matchedMetric, partiallyMatchedMetric, notMatchedMetric)
	}

	return result
}

func generateMatchedMetric(terms TaggedTermList) string {
	pathParts := make([]string, 0)
	for _, term := range terms {
		tagValues := strings.Split(term.Value, "|")
		tagValue := tagValues[rand.Intn(len(tagValues))]
		pathParts = addTag(pathParts, term, tagValue)
	}
	return strings.Join(pathParts, ";")
}

func generatePartiallyMatchedMetric(terms TaggedTermList) string {
	// there will be only one matched tag
	matchedTag := terms[rand.Intn(len(terms))]
	randomTagValueLength := 5

	pathParts := make([]string, 0, len(terms))
	for _, term := range terms {
		var tagValue string
		if term.Key == matchedTag.Key {
			tagValues := strings.Split(term.Value, "|")
			tagValue = tagValues[rand.Intn(len(tagValues))]
		} else {
			tagValue = randStringBytes(randomTagValueLength)
		}
		pathParts = addTag(pathParts, term, tagValue)
	}
	return strings.Join(pathParts, ";")
}

func generateNotMatchedMetric(terms TaggedTermList) string {
	randomTagValueLength := 5
	pathParts := make([]string, 0)
	for _, tagSpec := range terms {
		tagValue := randStringBytes(randomTagValueLength)
		pathParts = addTag(pathParts, tagSpec, tagValue)
	}
	return strings.Join(pathParts, ";")
}

func addTag(pathParts []string, term TaggedTerm, tagValue string) []string {
	if term.Key == "__name__" {
		// name tag should be prepended
		return append([]string{tagValue}, pathParts...)
	}
	return append(pathParts, term.Key+"="+tagValue)
}

const letterBytes = "abcdefghijklmnopqrstuvwxyz"

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
