package json

import (
	"fmt"
	"lkserver/internal/models"
	"lkserver/internal/repository"
)

type rankHistoryRepo struct {
	dataDir     string
	rankHistory []*models.RankHistory
}

func (r *repo) initRankHistoryRepo() error {
	repo := &rankHistoryRepo{dataDir: r.dataDir}
	if err := repo.init(r.individuals, r.ranks); err != nil {
		return err
	}

	r.ranksHistory = repo

	return nil
}

func (r *rankHistoryRepo) GetLast(iin string) (*models.RankHistory, error) {
	allRanks, err := r.GetHistory(iin)
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

func (r *rankHistoryRepo) GetHistory(iin string) ([]*models.RankHistory, error) {
	var ranks []*models.RankHistory
	for _, rank := range r.rankHistory {
		if rank.Individual.IndividualNumber == iin {
			ranks = append(ranks, rank)
		}
	}

	return ranks, nil

}

func (r *rankHistoryRepo) Close() {}

func (jr *rankHistoryRepo) init(i repository.IndividualsProvider, r repository.RankProvider) error {
	data := []struct {
		Date       models.JSONTime
		Rank       models.JSONByte
		Individual models.JSONByte
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
			if err == models.ErrNotFound {
				return models.ErrRefIntegrity
			}
			return err
		}
		rank, err := r.Get(row.Rank)
		if err != nil {
			if err == models.ErrNotFound {
				return models.ErrRefIntegrity
			}
			return err
		}
		history = append(history, &models.RankHistory{
			Date:       row.Date,
			Individual: individual,
			Rank:       rank,
		})
	}

	jr.rankHistory = history

	return nil
}
