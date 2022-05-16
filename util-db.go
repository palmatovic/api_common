package api_common

import (
	"fmt"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"

	log "github.com/sirupsen/logrus"
	gsql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// GetDB open new connection pool.
// This method have to be invokated only once, maybe you have to make some tuning for pool size
func GetDB(serviceConfig *MicroserviceConfiguration, logLevel logger.LogLevel) (*gorm.DB, error) {
	log.Traceln("calling GetDB method")
	if serviceConfig == nil {
		log.Errorln("cannot get database connection because service configuration is not initialized")
		os.Exit(EXIT_CODE_MISSING_CONFIG)
	}
	/*var typeTLS string
	if serviceConfig.Infrastructure.Database.SslEnabled {
		// skip-verify for self-signed
		typeTLS = "custom"
	} else {
		log.Errorln("cannot connect to the database in plain mode, only tls/ssl is allowed")
		return nil, fmt.Errorf("cannot connect to the database in plain mode, only tls/ssl connection is allowed")
	}*/
	dbPassword, errGetPasswordDB := GetSecretString(serviceConfig.Infrastructure.Database.PasswordFilepath)
	if errGetPasswordDB != nil {
		log.WithError(errGetPasswordDB).Errorf("cannot get database password secret: %s", serviceConfig.Infrastructure.Database.PasswordFilepath)
		return nil, fmt.Errorf("cannot get database password secret: %s", serviceConfig.Infrastructure.Database.PasswordFilepath)
	}
	configDB := mysql.Config{
		User:                 serviceConfig.Infrastructure.Database.Username,
		Passwd:               dbPassword,
		Addr:                 fmt.Sprintf("%s:%d", serviceConfig.Infrastructure.Database.Address, serviceConfig.Infrastructure.Database.Port),
		Net:                  "tcp",
		DBName:               serviceConfig.Infrastructure.Database.Name,
		Loc:                  time.UTC,
		ParseTime:            true,
		AllowNativePasswords: true,
		//TLSConfig:            typeTLS,
	}

	connectionString := configDB.FormatDSN()

	/*fileCA := serviceConfig.Infrastructure.Common.InternalCACertFilepath

	rootCertPool := x509.NewCertPool()
	CA, errReadCA := os.ReadFile(fileCA)
	if errReadCA != nil {
		log.WithError(errReadCA).Errorf("cannot read file ca: %s", fileCA)
		return nil, fmt.Errorf("cannot read file ca: %s", fileCA)
	}
	if validCA := rootCertPool.AppendCertsFromPEM(CA); !validCA {
		log.Errorf("failed to append ca from pem file: %s", fileCA)
		return nil, fmt.Errorf("failed to append ca from pem file: %s", fileCA)
	}

		err := mysql.RegisterTLSConfig(typeTLS, &tls.Config{
			RootCAs: rootCertPool,
		})
		if err != nil {
			return nil, err
		}

	*/
	DB, errOpenGORM := gorm.Open(gsql.Open(connectionString), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if errOpenGORM != nil {
		log.WithError(errOpenGORM).Error("cannot open gorm connection")
		return nil, fmt.Errorf("cannot open gorm connection")
	}
	return DB, nil
}
