package aquaticket

import "time"

type Distributeur struct {
	curseurPassage int
	curseurAttente int
}

func NouveauDistributeur() Distributeur {
	return Distributeur{curseurPassage: 1, curseurAttente: 0}
}

func (dist *Distributeur) NouveauTicket() int {
	dist.curseurAttente++
	return dist.curseurAttente
}

func (dist *Distributeur) PassageFini() {
	dist.curseurPassage++
}

func (dist *Distributeur) PeutPasser(ticket int) bool {
	return ticket == dist.curseurPassage
}

func (dist *Distributeur) ExecutionQuandTicketPret(fonction func() error) error {
	var numTicket int = dist.NouveauTicket()
	defer dist.PassageFini()
	var err error
	if dist.PeutPasser(numTicket) {
		err = fonction()
	} else {
		ticker := time.NewTicker(10 * time.Millisecond)
		for range ticker.C {
			if dist.PeutPasser(numTicket) {
				err = fonction()
				ticker.Stop()
				return err
			}
		}
		time.Sleep(30 * time.Second)
	}
	return err
}
