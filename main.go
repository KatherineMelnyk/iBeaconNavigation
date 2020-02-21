package main

import (
	"fmt"
	"math"
)

type beacon struct {
	name string
	x, y float64
	p0   float64
}

type obstacle struct {
	id, typeOfObstacle int
	x, y               float64
}

type object struct {
	lines []int
	x, y  float64
	dist  []float64
}

func distanceToObject(ob object, b beacon) float64 {
	return math.Sqrt(math.Pow(ob.x-b.x, 2.) + math.Pow(ob.y-b.y, 2.))
}

func distanceRSS(rss int) float64 {
	p0 := -80
	d0 := 1.
	n := 2.
	numerator := rss - p0
	denominator := -10 * n
	return math.Pow(10, float64(numerator)/denominator) * d0
}

func meanSquareError(ob object, beacons []beacon) float64 {
	distanceCalculated := 0.
	for i, _ := range beacons {
		distanceCalculated += distanceToObject(ob, beacons[i])
	}
	mse := math.Pow(distanceCalculated, 2.)
	return mse / float64(len(beacons))
}

func main() {
	beacons := beaconsValue("geoBeacons.csv")
	fmt.Printf("Your beacon is %s with coordinates (%.4f,%.4f) \n",
		beacons[0].name, beacons[0].x, beacons[0].y)
	//obstacles := obstaclesValue("geoObstacles.csv")
	//fmt.Printf("Your obstacle id is %v with coordinates (%.4f,%.4f) "+
	//	"and type %v \n", obstacles[0].id, obstacles[0].x,
	//	obstacles[0].y, obstacles[0].typeOfObstacle)
	rss := RSSIMeasurements("geo_rss.csv")
	fmt.Println(rss[0])
	objs := coordinates("geo_crd.csv")
	fmt.Println(objs[0])

	for i, obj := range objs {
		for _, b := range beacons {
			obj.dist = append(obj.dist, distanceToObject(obj, b))
		}
		objs[i] = obj
	}
	fmt.Println(objs[0])
	//objs[0].CalculateCalibrationPower(rss)
	CalculateCalibrationPowers(beacons, objs, rss)
	var err float64
	for _, obj := range objs {
		_, _, e := obj.Coordinate(rss, beacons)
		err += e
	}
	fmt.Println(err / float64(len(objs)))
}

func CalculateCalibrationPowers(beacons []beacon, objs []object, rss [][]int) {
	objInd := map[int]int{} // line index -> object index
	for i, obj := range objs {
		for _, j := range obj.lines {
			objInd[j] = i
		}
	}
	for i, b := range beacons {
		var cnt float64
		for j, vals := range rss {
			rssi := vals[i]
			if rssi == 100 {
				continue
			}
			obj := objs[objInd[j]]
			b.p0 += float64(rssi) + 10*2*math.Log10(obj.dist[i])
			cnt++
		}
		b.p0 /= cnt
		//fmt.Println(b.p0)
		beacons[i] = b
	}
}

func (b *beacon) dist(rssi float64) float64 {
	return math.Pow(10, (rssi-b.p0)/-10/2)
}

func (obj *object) Coordinate(
	rss [][]int, beacons []beacon,
) (float64, float64, float64) {
	var d []float64
	var x []float64
	var y []float64

	for _, i := range obj.lines {
		for j, rssi := range rss[i] {
			if rssi == 100 {
				continue
			}
			b := beacons[j]
			d = append(d, b.dist(float64(rssi)))
			//fmt.Println(b.dist(float64(rssi)))
			x = append(x, b.x)
			y = append(y, b.y)
		}
	}

	fmt.Println(len(d))

	var xc, yc float64
	for i := 0; i < 10000; i++ {
		var dfdx float64
		var dfdy float64
		for i, xi := range x {
			yi := y[i]
			di := d[i]

			k := 1 / math.Sqrt((xc-xi)*(xc-xi)+(yc-yi)*(yc-yi))
			dfdx += (xc - xi) * 2 * (1 - di*k)
			dfdy += (yc - yi) * 2 * (1 - di*k)
		}
		dd := math.Sqrt(dfdx*dfdx + dfdy*dfdy)
		xc -= 0.01 * dfdx / dd
		yc -= 0.01 * dfdy / dd
		//fmt.Printf("(%7.3f %7.3f)\n", xc, yc)
	}
	err := math.Sqrt((xc-obj.x)*(xc-obj.x) + (yc-obj.y)*(yc-obj.y))
	fmt.Printf("(%7.3f %7.3f) (%7.3f %7.3f) %7.3f\n",
		xc, yc, obj.x, obj.y, err)
	return xc, yc, err
}
