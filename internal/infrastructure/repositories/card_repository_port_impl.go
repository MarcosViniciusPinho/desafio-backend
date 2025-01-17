package repositories

import (
	"account/internal/domain"
	"account/internal/domain/ports/outbounds"
	"account/internal/infrastructure/repositories/config"
	"account/internal/infrastructure/repositories/entities"
	"account/pkg/log"
)

type CardRepositoryPortImpl struct{}

func (c *CardRepositoryPortImpl) Create(personId int, card domain.Card) *domain.Card {
	conn := config.OpenConnection()
	defer conn.Close()

	cardEntity := entities.NewCard(
		card.Title,
		card.CardNumber(),
		card.ExpireMonth,
		card.ExpireYear,
		card.EncryptSecurityCode(),
		personId,
	)

	log.Info(
		"entity",
		cardEntity,
		"Information before being persisted",
	)

	_, err := conn.Model(cardEntity).Insert()

	if err != nil {
		panic(err)
	}

	log.Info(
		"id",
		cardEntity.Id,
		"Id that was generated after recording in the database",
	)

	return domain.NewCardFull(
		cardEntity.Id,
		cardEntity.Title,
		cardEntity.Pan,
		cardEntity.ExpireMonth,
		cardEntity.ExpireYear,
		cardEntity.SecurityCode,
		cardEntity.Date,
	)
}

func (c *CardRepositoryPortImpl) Update(personId, id int, card domain.Card) *domain.Card {
	conn := config.OpenConnection()
	defer conn.Close()

	cardEntity := entities.NewCard(
		card.Title,
		card.CardNumber(),
		card.ExpireMonth,
		card.ExpireYear,
		card.EncryptSecurityCode(),
		personId,
	)

	log.Info(
		"entity",
		cardEntity,
		"Information before being persisted",
	)

	query := conn.Model(cardEntity).Where("id = ?", id)

	_, err := query.Update()

	if err != nil {
		panic(err)
	}

	log.InfoSimple("The data has been updated correctly")

	return domain.NewCardFull(
		cardEntity.Id,
		cardEntity.Title,
		cardEntity.Pan,
		cardEntity.ExpireMonth,
		cardEntity.ExpireYear,
		cardEntity.SecurityCode,
		cardEntity.Date,
	)
}

func (c *CardRepositoryPortImpl) ExistsByPersonIdAndId(personId, id int) bool {
	conn := config.OpenConnection()
	defer conn.Close()

	cardEntity := entities.NewCardDefault()
	query := conn.Model(cardEntity).
		Where("id = ?", id).
		Where("people_id = ?", personId)

	log.Info(
		"query",
		"where id = ? and people_id = ?",
		"Searching by id and people_id",
	)

	foundCard, err := query.Exists()

	if err != nil {
		panic(err)
	}
	return foundCard
}

func (c *CardRepositoryPortImpl) FindById(id int) *domain.Card {
	conn := config.OpenConnection()
	defer conn.Close()

	var cardEntity entities.Card
	query := conn.Model(&cardEntity).Where("id = ?", id)

	log.Info(
		"query",
		"where id = ?",
		"Searching by id",
	)

	foundCard, err := query.Exists()

	if err != nil {
		panic(err)
	}
	if !foundCard {
		return nil
	}

	err = query.Select()
	if err != nil {
		panic(err)
	}

	return domain.NewCardFull(
		cardEntity.Id,
		cardEntity.Title,
		cardEntity.Pan,
		cardEntity.ExpireMonth,
		cardEntity.ExpireYear,
		cardEntity.SecurityCode,
		cardEntity.Date,
	)
}

func (c *CardRepositoryPortImpl) Delete(id int) {
	conn := config.OpenConnection()
	defer conn.Close()

	var cardEntity entities.Card
	conn.Model(&cardEntity).Where("id = ?", id).Delete()

	log.InfoSimple("Delete completed successfully")
}

func (c *CardRepositoryPortImpl) FindAllByPersonId(personId int) []domain.Card {
	conn := config.OpenConnection()
	defer conn.Close()

	log.Info(
		"query",
		"where people_id = ?",
		"Searching by people_id",
	)

	var cardEntities []entities.Card
	err := conn.Model(&cardEntities).Where("people_id = ?", personId).Select()
	if err != nil {
		panic(err)
	}

	var cardDomains []domain.Card
	for _, cardEntity := range cardEntities {
		cardDomains = append(
			cardDomains,
			domain.NewCard(
				cardEntity.Id,
				cardEntity.Title,
				cardEntity.Pan,
				cardEntity.ExpireMonth,
				cardEntity.ExpireYear,
				cardEntity.SecurityCode,
				cardEntity.Date,
			),
		)
	}
	return cardDomains
}

func NewCardRepositoryPort() outbounds.CardRepositoryPort {
	return &CardRepositoryPortImpl{}
}
