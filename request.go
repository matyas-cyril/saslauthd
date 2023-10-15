package saslauthd

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/matyas-cyril/logme"
)

func requestClient(cnx net.Conn, conf configFile, msgID logme.MsgID) error {

	// Context pour le TimeOut
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(conf.Server.ClientTimeout)*time.Second)
	defer func() {

		if Debug() {
			debug.addLogInFile(fmt.Sprintf("#[%s] -> .. -> close client cnx", msgID))
		}

		clients.Done(1)
		cancel()

	}()

	// Déclaration des channels en cas de succès (OK ou) et de panic de go routine
	done := make(chan request, 1)
	panicChan := make(chan interface{}, 1)

	// Go routine d'exécution du traitement métier du client
	go func() {

		defer func() {
			if p := recover(); p != nil {
				panicChan <- p
			}
		}()

		done <- handleConnection(cnx, conf, msgID)
	}()

	select {

	// Traitement d'un retour normal
	case req := <-done:

		if Debug() {
			debug.addLogInFile(fmt.Sprintf("#[%s] -> .. -> write socket: %v", msgID, req.auth))
		}

		if req.auth {

			if conf.Server.Stat > 0 {
				statClientOK.Inc()
			}
			cnx.Write([]byte(SASL_SUCCESS))
		} else {
			if conf.Server.Stat > 0 {
				statClientKO.Inc()
			}
			cnx.Write([]byte(SASL_FAIL))
		}

		if req.err != nil {
			return req.err
		}
		return nil

	// Pb de la goroutine
	case p := <-panicChan:

		if Debug() {
			debug.addLogInFile(fmt.Sprintf("#[%s] -> .. -> write socket: false", msgID))
		}

		if conf.Server.Stat > 0 {
			statClientKO.Inc()
		}

		cnx.Write([]byte(SASL_FAIL))
		panic(p)

	// TimeOut
	case <-ctx.Done():

		if Debug() {
			debug.addLogInFile(fmt.Sprintf("#[%s] -> .. -> client timeout cnx", msgID))
			debug.addLogInFile(fmt.Sprintf("#[%s] -> .. -> write socket: false", msgID))
		}

		if conf.Server.Stat > 0 {
			statClientKO.Inc()
		}

		cnx.Write([]byte(SASL_FAIL))
		close(done)

		return errors.New("TimeOut")

	}

}
