package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"
)

type Stat struct {
	Ackr float64 `json:"ackr,omitempty"` // Percentage of upstream datagrams that were acknowledged
	Alti int     `json:"alti,omitempty"` // GPS altitude of the gateway in meter RX (integer)
	Dwnb uint    `json:"dwnb,omitempty"` // Number of downlink datagrams received (unsigned integer)
	Lati float64 `json:"lati,omitempty"` // GPS latitude of the gateway in degree (float, N is +)
	Long float64 `json:"long,omitempty"` // GPS latitude of the gateway in dgree (float, E is +)
	Rxfw uint    `json:"rxfw,omitempty"` // Number of radio packets forwarded (unsigned integer)
	Rxnb uint    `json:"rxnb,omitempty"` // Number of radio packets received (unsigned integer)
	Rxok uint    `json:"rxok,omitempty"` // Number of radio packets received with a valid PHY CRC
	Time *mytime  `json:"time,omitempty"` // UTC 'system' time of the gateway, ISO 8601 'expanded' format
	Txnb uint    `json:"txnb,omitempty"` // Number of packets emitted (unsigned integer)
}

// RXPK represents an uplink json message format sent by the gateway
type RXPK struct {
	Chan uint    `json:"chan,omitempty"` // Concentrator "IF" channel used for RX (unsigned integer)
	Codr string  `json:"codr,omitempty"` // LoRa ECC coding rate identifier
	Data string  `json:"data,omitempty"` // Base64 encoded RF packet payload, padded
	Datr *datr    `json:"datr,omitempty"` // FSK datarate (unsigned in bit per second) || LoRa datarate identifier
	Freq float64 `json:"freq,omitempty"` // RX Central frequency in MHx (unsigned float, Hz precision)
	Lsnr float64 `json:"lsnr,omitempty"` // LoRa SNR ratio in dB (signed float, 0.1 dB precision)
	Modu string  `json:"modu,omitempty"` // Modulation identifier "LORA" or "FSK"
	Rfch uint    `json:"rfch,omitempty"` // Concentrator "RF chain" used for RX (unsigned integer)
	Rssi int     `json:"rssi,omitempty"` // RSSI in dBm (signed integer, 1 dB precision)
	Size uint    `json:"size,omitempty"` // RF packet payload size in bytes (unsigned integer)
	Stat int     `json:"stat,omitempty"` // CRC status: 1 - OK, -1 = fail, 0 = no CRC
	Time *mytime  `json:"time,omitempty"` // UTC time of pkt RX, us precision, ISO 8601 'compact' format
	Tmst uint    `json:"tmst,omitempty"` // Internal timestamp of "RX finished" event (32b unsigned)
}

// TXPK represents a downlink json message format received by the gateway.
// Most field are optional.
type TXPK struct {
	Codr string  `json:"codr,omitempty"` // LoRa ECC coding rate identifier
	Data string  `json:"data,omitempty"` // Base64 encoded RF packet payload, padding optional
	Datr *datr    `json:"datr,omitempty"` // LoRa datarate identifier (eg. SF12BW500) || FSK Datarate (unsigned, in bits per second)
	Fdev uint    `json:"fdev,omitempty"` // FSK frequency deviation (unsigned integer, in Hz)
	Freq float64 `json:"freq,omitempty"` // TX central frequency in MHz (unsigned float, Hz precision)
	Imme bool    `json:"imme,omitempty"` // Send packet immediately (will ignore tmst & time)
	Ipol bool    `json:"ipol,omitempty"` // Lora modulation polarization inversion
	Modu string  `json:"modu,omitempty"` // Modulation identifier "LORA" or "FSK"
	Ncrc bool    `json:"ncrc,omitempty"` // If true, disable the CRC of the physical layer (optional)
	Powe uint    `json:"powe,omitempty"` // TX output power in dBm (unsigned integer, dBm precision)
	Prea uint    `json:"prea,omitempty"` // RF preamble size (unsigned integer)
	Rfch uint    `json:"rfch,omitempty"` // Concentrator "RF chain" used for TX (unsigned integer)
	Size uint    `json:"size,omitempty"` // RF packet payload size in bytes (unsigned integer)
	Time *mytime  `json:"time,omitempty"` // Send packet at a certain time (GPS synchronization required)
	Tmst uint    `json:"tmst,omitempty"` // Send packet on a certain timestamp value (will ignore time)
}

type datr struct {
	kind  string
	value string
}

func (d *datr) MarshalJSON() ([]byte, error) {
	if d.kind == "uint" {
		return []byte(d.value), nil
	}
	return append(append([]byte(`"`), []byte(d.value)...), []byte(`"`)...), nil
}

type mytime struct {
	layout string
	value  time.Time
}

func (m *mytime) MarshalJSON() ([]byte, error) {
	return append(append([]byte(`"`), []byte(m.value.Format(m.layout))...), []byte(`"`)...), nil
}

type Payload struct {
	RXPK *[]RXPK `json:"rxpk,omitempty"`
	Stat *Stat   `json:"stat,omitempty"`
	TXPK *TXPK   `json:"txpk,omitempty"`
}

func main() {
	time1, _ := time.Parse(time.RFC3339Nano, "2013-03-31T16:21:17.528002Z")
	time2, _ := time.Parse(time.RFC3339Nano, "2013-03-31T16:21:17.530974Z")
	time3, _ := time.Parse(time.RFC3339, "2014-01-12T08:59:28Z")

	rawRXPKs, _ := json.Marshal(Payload{
		RXPK: &[]RXPK{
			RXPK{
				Time: &mytime{time.RFC3339Nano, time1},
				Tmst: 3512348611,
				Chan: 2,
				Rfch: 0,
				Freq: 866.349812,
				Stat: 1,
				Modu: "LORA",
				Datr: &datr{"string", "SF7BW125"},
				Codr: "4/6",
				Rssi: -35,
				Lsnr: 5.1,
				Size: 32,
				Data: "-DS4CGaDCdG+48eJNM3Vai-zDpsR71Pn9CPA9uCON84",
			},
			RXPK{
				Chan: 9,
				Data: "VEVTVF9QQUNLRVRfMTIzNA==",
				Datr: &datr{"uint", "50000"},
				Freq: 869.1,
				Modu: "FSK",
				Rfch: 1,
				Rssi: -75,
				Size: 16,
				Stat: 1,
				Time: &mytime{time.RFC3339Nano, time2},
				Tmst: 3512348514,
			},
		},
	})

	rawStat, _ := json.Marshal(Payload{
		Stat: &Stat{
			Ackr: 100.0,
			Alti: 145,
			Long: 3.25230,
			Rxok: 2,
			Rxfw: 2,
			Rxnb: 2,
			Lati: 46.24,
			Dwnb: 2,
			Txnb: 2,
			Time: &mytime{"2006-01-02 15:04:05 GMT", time3},
		},
	})

	rawRXPKsStat, _ := json.Marshal(Payload{
		RXPK: &[]RXPK{
			RXPK{
				Time: &mytime{time.RFC3339Nano, time1},
				Tmst: 3512348611,
				Chan: 2,
				Rfch: 0,
				Freq: 866.349812,
				Stat: 1,
				Modu: "LORA",
				Datr: &datr{"string", "SF7BW125"},
				Codr: "4/6",
				Rssi: -35,
				Lsnr: 5.1,
				Size: 32,
				Data: "-DS4CGaDCdG+48eJNM3Vai-zDpsR71Pn9CPA9uCON84",
			},
			RXPK{
				Chan: 9,
				Data: "VEVTVF9QQUNLRVRfMTIzNA==",
				Datr: &datr{"uint", "50000"},
				Freq: 869.1,
				Modu: "FSK",
				Rfch: 1,
				Rssi: -75,
				Size: 16,
				Stat: 1,
				Time: &mytime{time.RFC3339Nano, time2},
				Tmst: 3512348514,
			},
		},
		Stat: &Stat{
			Ackr: 100.0,
			Alti: 145,
			Long: 3.25230,
			Rxok: 2,
			Rxfw: 2,
			Rxnb: 2,
			Lati: 46.24,
			Dwnb: 2,
			Txnb: 2,
			Time: &mytime{"2006-01-02 15:04:05 GMT", time3},
		},
	})

	rawTXPK, _ := json.Marshal(Payload{
		TXPK: &TXPK{
			Imme: true,
			Freq: 864.123456,
			Rfch: 0,
			Powe: 14,
			Modu: "LORA",
			Datr: &datr{"string", "SF11BW125"},
			Codr: "4/6",
			Ipol: false,
			Size: 32,
			Data: "H3P3N2i9qc4yt7rK7ldqoeCVJGBybzPY5h1Dd7P7p8v",
		},
	})

	fmt.Printf("Stat:       %v\n", string(rawStat))
	fmt.Printf("Raw Stat:      ")
	for _, x := range rawStat {
		fmt.Printf("0x%x,", x)
	}
	fmt.Printf("\n\n\n")
	ioutil.WriteFile("./test_data/marshal_stat", rawStat, 0644)

	fmt.Printf("RXPKs:       %v\n", string(rawRXPKs))
	fmt.Printf("Raw RXPKs:      ")
	for _, x := range rawRXPKs {
		fmt.Printf("0x%x,", x)
	}
	fmt.Printf("\n\n\n")
	ioutil.WriteFile("./test_data/marshal_rxpk", rawRXPKs, 0644)

	fmt.Printf("RXPKsStats:       %v\n", string(rawRXPKsStat))
	fmt.Printf("Raw RXPKsStats:      ")
	for _, x := range rawRXPKsStat {
		fmt.Printf("0x%x,", x)
	}
	fmt.Printf("\n\n\n")
	ioutil.WriteFile("./test_data/marshal_rxpk_stat", rawRXPKsStat, 0644)

	fmt.Printf("TXPK:       %v\n", string(rawTXPK))
	fmt.Printf("Raw TXPK:      ")
	for _, x := range rawTXPK {
		fmt.Printf("0x%x,", x)
	}
	fmt.Printf("\n\n\n")
	ioutil.WriteFile("./test_data/marshal_txpk", rawTXPK, 0644)
}
