package services

import (
	"errors"
	"fmt"
	"time"

	"tcc-tech/queue-backend/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrQueueExhausted = errors.New("all queue numbers exhausted (A0-Z9), please clear the queue")
	bangkokLoc        = time.FixedZone("Asia/Bangkok", 7*60*60)
)

type QueueService struct {
	db *gorm.DB
}

func NewQueueService(db *gorm.DB) *QueueService {
	return &QueueService{db: db}
}

// IssueNextTicket ออกบัตรคิวใหม่ พร้อม concurrency control ด้วย SELECT FOR UPDATE
func (s *QueueService) IssueNextTicket() (*models.QueueTicket, error) {
	var ticket models.QueueTicket

	err := s.db.Transaction(func(tx *gorm.DB) error {
		// 1. Lock singleton row
		var state models.QueueState
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&state, 1).Error; err != nil {
			return err
		}

		// 2. Compute next queue number
		nextLetter, nextNumber, err := computeNext(state.CurrentLetter, state.CurrentNumber)
		if err != nil {
			return err
		}

		// 3. Update state
		state.CurrentLetter = nextLetter
		state.CurrentNumber = nextNumber
		if err := tx.Save(&state).Error; err != nil {
			return err
		}

		// 4. Insert ticket
		ticket = models.QueueTicket{
			QueueNumber: fmt.Sprintf("%s%d", nextLetter, nextNumber),
			IssuedAt:    time.Now().In(bangkokLoc),
		}
		return tx.Create(&ticket).Error
	})

	if err != nil {
		return nil, err
	}
	return &ticket, nil
}

// GetCurrentQueue ดึงคิวปัจจุบัน
func (s *QueueService) GetCurrentQueue() (string, *time.Time, error) {
	var state models.QueueState
	if err := s.db.First(&state, 1).Error; err != nil {
		return "00", nil, err
	}

	// Queue is cleared
	if state.CurrentLetter == "" {
		return "00", nil, nil
	}

	queueNumber := fmt.Sprintf("%s%d", state.CurrentLetter, state.CurrentNumber)

	// Get latest ticket for timestamp
	var ticket models.QueueTicket
	if err := s.db.Where("queue_number = ?", queueNumber).
		Order("issued_at DESC").First(&ticket).Error; err != nil {
		return queueNumber, nil, nil
	}

	return queueNumber, &ticket.IssuedAt, nil
}

// ClearQueue ล้างคิวทั้งหมด reset เป็น "00"
func (s *QueueService) ClearQueue() error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 1. Reset state
		if err := tx.Model(&models.QueueState{}).Where("id = ?", 1).
			Updates(map[string]interface{}{
				"current_letter": "",
				"current_number": -1,
			}).Error; err != nil {
			return err
		}

		// 2. Delete all tickets
		return tx.Where("1 = 1").Delete(&models.QueueTicket{}).Error
	})
}

// computeNext คำนวณหมายเลขคิวถัดไป
// "" / -1 (cleared) → A / 0
// A-Y / 9          → next letter / 0
// Z / 9            → error (exhausted)
// any / 0-8        → same letter / number+1
func computeNext(letter string, number int) (string, int, error) {
	// Cleared state → start from A0
	if letter == "" {
		return "A", 0, nil
	}

	// Check exhausted
	if letter == "Z" && number == 9 {
		return "", 0, ErrQueueExhausted
	}

	// Number < 9 → increment number
	if number < 9 {
		return letter, number + 1, nil
	}

	// Number == 9 → next letter, reset number to 0
	nextLetter := string(rune(letter[0]) + 1)
	return nextLetter, 0, nil
}
