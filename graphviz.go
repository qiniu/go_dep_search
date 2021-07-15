// +build graphviz

package main

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
)

func resultToSvg(result map[string][]string) {
	g := graphviz.New()
	defer g.Close()
	graph, err := g.Graph()
	if err != nil {
		log.Panicln(err)
	}
	defer graph.Close()
	graph.SetRankDir(cgraph.LRRank)
	for from, tos := range result {
		fromNode, err := graph.CreateNode(from)
		if err != nil {
			log.Panicln(err)
		}
		fromNode.SetStyle(cgraph.BoldNodeStyle)
		for _, to := range tos {
			toNode, err := graph.CreateNode(to)
			if err != nil {
				log.Panicln(err)
			}
			toNode.SetStyle(cgraph.BoldNodeStyle)
			edge, err := graph.CreateEdge(from, fromNode, toNode)
			if err != nil {
				log.Panicln(err)
			}
			edge.SetURL("#" + from + "->" + to)
			edge.SetStyle(cgraph.BoldEdgeStyle)
		}
	}
	format := graphviz.Format(strings.Trim(filepath.Ext(*graphResultFile), "."))
	err = g.RenderFilename(graph, format, *graphResultFile)
	if err != nil {
		log.Panicln(err)
	}
	fmt.Println("result saved to " + *graphResultFile)
}
