package main

import (
	"bufio"
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
)

func BeaconsValue(file string) []Beacon {
	csvFile, _ := os.Open(file)
	reader := csv.NewReader(bufio.NewReader(csvFile))
	var beacons []Beacon
	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}
		X, _ := strconv.ParseFloat(line[1], 64)
		Y, _ := strconv.ParseFloat(line[2], 64)
		beacons = append(beacons, Beacon{
			name: line[0],
			x:    X,
			y:    Y,
		})
	}
	return beacons[1:]
}

func obstaclesValue(file string) []Obstacle {
	csvFile, _ := os.Open(file)
	reader := csv.NewReader(bufio.NewReader(csvFile))
	var obstacles []Obstacle
	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}
		ID, _ := strconv.ParseFloat(line[0], 64)
		X, _ := strconv.ParseFloat(line[1], 64)
		Y, _ := strconv.ParseFloat(line[2], 64)
		typeObstacle, _ := strconv.ParseFloat(line[3], 64)
		obstacles = append(obstacles, Obstacle{
			id:             int(ID),
			x:              X,
			y:              Y,
			typeOfObstacle: int(typeObstacle),
		})
	}
	return obstacles[1:]
}

func RSSIMeasurements(file string) [][]int {
	csvFile, _ := os.Open(file)
	reader := csv.NewReader(bufio.NewReader(csvFile))
	var rss [][]int
	for {
		var lineRss []int
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		for i := 0; i < len(line); i++ {
			signal, _ := strconv.ParseFloat(line[i], 64)
			lineRss = append(lineRss, int(signal))
		}
		rss = append(rss, lineRss)
	}
	return rss
}

func ObjectValue(file string) []Object {
	csvFile, _ := os.Open(file)
	reader := csv.NewReader(bufio.NewReader(csvFile))
	var objects []Object
	var obj Object
	for i := 0; ; i++ {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		if len(line) != 2 {
			log.Fatal("line length should be 2")
		}
		x, _ := strconv.ParseFloat(line[0], 64)
		y, _ := strconv.ParseFloat(line[1], 64)
		if obj.x == x && obj.y == y {
			obj.lines = append(obj.lines, i)
		} else {
			objects = append(objects, obj)
			obj = Object{lines: []int{i}, x: x, y: y}
		}
	}
	objects = append(objects, obj)
	return objects[1:]
}
