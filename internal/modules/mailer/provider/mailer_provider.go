package provider

import (
	"github.com/RakaMurdiarta/online-shop-system/internal/modules/mailer/services"
	mailerServiceImpl "github.com/RakaMurdiarta/online-shop-system/internal/modules/mailer/services/Impl"
	"github.com/RakaMurdiarta/online-shop-system/pkg/mailer"
)

// MailerProvider wires the MailerService against any mailer.Transport
// implementation, keeping the service decoupled from the concrete provider.
func MailerProvider(transport mailer.Transport) services.MailerService {
	return mailerServiceImpl.NewMailerService(transport)
}
