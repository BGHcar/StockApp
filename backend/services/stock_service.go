package services

import (
	"backend/models"
	"backend/repositories"
	"fmt"
	"log"
	"sort"
	"strings"
)

func isPositiveRating(rating string) bool {
	ratingLower := strings.ToLower(rating)
	positiveRatings := []string{"buy", "strong buy", "outperform", "overweight", "accumulate"}
	for _, pr := range positiveRatings {
		if strings.Contains(ratingLower, pr) {
			return true
		}
	}
	return false
}

func isNegativeRating(rating string) bool {
	ratingLower := strings.ToLower(rating)
	negativeRatings := []string{"sell", "strong sell", "underperform", "underweight", "reduce"}
	for _, nr := range negativeRatings {
		if strings.Contains(ratingLower, nr) {
			return true
		}
	}
	return false
}

func isUpgrade(from, to string) bool {
	// Simple logic: Upgrade if 'to' is positive and 'from' is not positive (or is negative/neutral)
	return isPositiveRating(to) && !isPositiveRating(from)
}

func GetRecommendationsService() ([]models.Recommendation, error) {

	scores := make(map[string]*models.Recommendation)
	reasons := make(map[string][]string) // Store reasons per ticker

	stocks, err := repositories.GetByRecommendation()
	if err != nil {
		log.Println("Error fetching recommendations:", err)
		return nil, fmt.Errorf("error fetching recommendations: %v", err)
	}

	for _, stock := range stocks {
		if stock.Ticker == "" {
			continue // Skip if Ticker is empty
		}

		currentScore := 0.0
		reasonFragments := []string{}

		// Score based on Action
		if strings.EqualFold(stock.Action, "Buy") {
			currentScore += 1.0
			reasonFragments = append(reasonFragments, "Buy Action")
		}

		// Score based on Rating To
		if isPositiveRating(stock.RatingTo) {
			currentScore += 1.0
			reasonFragments = append(reasonFragments, fmt.Sprintf("Positive Rating (%s)", stock.RatingTo))
		}

		// Score based on Upgrade
		if isUpgrade(stock.RatingFrom, stock.RatingTo) {
			currentScore += 0.5 // Bonus for upgrade
			reasonFragments = append(reasonFragments, fmt.Sprintf("Upgrade (%s -> %s)", stock.RatingFrom, stock.RatingTo))
		}

		// --- Potential Enhancements ---
		// - Weight by brokerage reputation
		// - Weight by recency (newer data gets higher score)
		// - Factor in target price upside (requires current price data)

		if currentScore > 0 {
			ticker := stock.Ticker
			if existing, ok := scores[ticker]; ok {
				// Update existing score and latest timestamp
				existing.Score += currentScore
				if stock.Time.After(existing.LastUpdate) {
					existing.LastUpdate = stock.Time
					existing.Company = stock.Company // Update company name too
				}
				reasons[ticker] = append(reasons[ticker], reasonFragments...)
			} else {
				// Create new recommendation entry
				scores[ticker] = &models.Recommendation{
					Ticker:     ticker,
					Company:    stock.Company,
					Score:      currentScore,
					LastUpdate: stock.Time,
					// Reason will be built later
				}
				reasons[ticker] = reasonFragments
			}
		}
	}

	recommendations := make([]models.Recommendation, 0, len(scores))
	for ticker, rec := range scores {
		uniqueReasons := make(map[string]bool)
		for _, reason := range reasons[ticker] {
			uniqueReasons[reason] = true
		}
		var reasonParts []string
		for part := range uniqueReasons {
			reasonParts = append(reasonParts, part)
		}
		sort.Strings(reasonParts)
		rec.Reason = strings.Join(reasonParts, ", ")

		recommendations = append(recommendations, *rec)
	}

	// Sort recommendations by score (highest first) :5

	sort.Slice(recommendations, func(i, j int) bool {
		return recommendations[i].Score > recommendations[j].Score
	})

	if limit := 5; len(recommendations) > limit {
		recommendations = recommendations[:limit]
	}

	log.Println(len(recommendations), " Recommendations generated successfully", recommendations)
	return recommendations, nil
}
