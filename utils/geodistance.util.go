package utils

import (
	"math"
	"sort"
	"zuck-my-clothe/zuck-my-clothe-backend/model"
)

// HaversineDistance calculates the distance between two coordinates in kilometers
func haversineDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371 // Radius of Earth in kilometers
	dLat := (lat2 - lat1) * math.Pi / 180
	dLon := (lon2 - lon1) * math.Pi / 180
	lat1 = lat1 * math.Pi / 180
	lat2 = lat2 * math.Pi / 180

	a := math.Sin(dLat/2)*math.Sin(dLat/2) + math.Cos(lat1)*math.Cos(lat2)*math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c
}

// BoundingBoxFilter filters branches within a bounding box
func boundingBoxFilter(userLat, userLon, radiusKm float64, branches []model.BranchDetail) []model.BranchDetail {
	// Convert radius to degrees
	latRange := radiusKm / 111.32
	lonRange := radiusKm / (111.32 * math.Cos(userLat*math.Pi/180))

	filteredBranches := []model.BranchDetail{}
	for _, branch := range branches {
		if math.Abs(branch.BranchLat-userLat) <= latRange && math.Abs(branch.BranchLon-userLon) <= lonRange {
			filteredBranches = append(filteredBranches, branch)
		}
	}
	return filteredBranches
}

// SortBranchesByDistance sorts branches by their distance from the user's location in ascending order
func SortBranchesByDistance(searchRadiusKm, userLat, userLon float64, branches []model.BranchDetail) []model.BranchDetail {
	filterdBranches := boundingBoxFilter(userLat, userLon, searchRadiusKm, branches)

	// Calculate distance for each shop
	for i := range filterdBranches {
		filterdBranches[i].Distance = haversineDistance(userLat, userLon, filterdBranches[i].BranchLat, filterdBranches[i].BranchLon)
	}

	// Sort branches by distance in ascending order
	sort.Slice(filterdBranches, func(i, j int) bool {
		return filterdBranches[i].Distance < filterdBranches[j].Distance
	})

	return filterdBranches
}
