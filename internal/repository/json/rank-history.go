package json

import (
	"fmt"
	"lkserver/internal/models"
	"lkserver/internal/repository"
)

type RankHistoryRepo struct {
	dataDir     string
	rankHistory []*models.RankHistory
}

func (r *repo) initRankHistoryRepo() error {
	repo := &RankHistoryRepo{dataDir: r.dataDir}
	if err := repo.init(r.individuals, r.ranks); err != nil {
		return err
	}

	r.ranksHistory = repo

	return nil
}

func (r *RankHistoryRepo) GetLast(people *models.Individuals) (*models.RankHistory, error) {
	allRanks, err := r.GetHistory(people)
	if err != nil {
		return nil, err
	}

	if len(allRanks) == 0 {
		return nil, fmt.Errorf("NOT FOUND")
	}

	result := allRanks[0]
	for _, rank := range allRanks[1:] {
		if rank.Date.After(result.Date) {
			result = rank
		}
	}

	return result, nil
}

func (r *RankHistoryRepo) GetHistory(people *models.Individuals) ([]*models.RankHistory, error) {
	var ranks []*models.RankHistory
	for _, rank := range r.rankHistory {
		if rank.Individual.IndividualNumber == people.IndividualNumber {
			ranks = append(ranks, rank)
		}
	}

	return ranks, nil

}

func (r *RankHistoryRepo) Close() {}

func (jr *RankHistoryRepo) init(i repository.IndividualsProvider, r repository.RankProvider) error {
	data := []struct {
		Date       models.JSONTime
		Rank       int
		Individual string
	}{}

	if err := initFile(
		fmt.Sprintf("%s/rank-history.json", jr.dataDir),
		&data,
	); err != nil {
		return err
	}
	var history []*models.RankHistory
	for _, row := range data {
		individual, err := i.Get(row.Individual)
		if err != nil {
			if err == repository.ErrNotFound {
				return repository.ErrRefIntegrity
			}
			return err
		}
		rank, err := r.Get(row.Rank)
		if err != nil {
			if err == repository.ErrNotFound {
				return repository.ErrRefIntegrity
			}
			return err
		}
		history = append(history, &models.RankHistory{
			Date:       row.Date,
			Individual: *individual,
			Rank:       *rank,
		})
	}

	jr.rankHistory = history

	return nil
}
