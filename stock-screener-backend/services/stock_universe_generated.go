package services

// This file is generated from market index sources.
// Sources:
// - TA universe: ../TASE_scanner/data/indexcomponents.csv
// - S&P 500: https://en.wikipedia.org/wiki/List_of_S%26P_500_companies
// - Nasdaq-100: https://en.wikipedia.org/wiki/Nasdaq-100

var taIndexSymbols = []string{
	"TSEM.TA", "ESLT.TA", "NVMI.TA", "POLI.TA", "LUMI.TA", "TEVA.TA", "DSCT.TA", "PHOE.TA", "MZTF.TA", "ENLT.TA",
	"NXSN.TA", "ORA.TA", "NICE.TA", "HARL.TA", "BEZQ.TA", "OPCE.TA", "AZRG.TA", "CLIS.TA", "CAMT.TA", "NVPT.TA",
	"FIBI.TA", "BIG.TA", "ICL.TA", "TASE.TA", "NWMD.TA", "MMHD.TA", "MGDL.TA", "MLSR.TA", "MGOR.TA", "DLEKG.TA",
	"SAE.TA", "PAZ.TA", "MVNE.TA", "STRS.TA", "SKBN.TA", "BSEN.TA", "ALHE.TA", "ENRG.TA", "PTNR.TA", "DORL.TA",
	"ENOG.TA", "AMOT.TA", "ELAL.TA", "RIT1.TA", "SPEN.TA", "MTAV.TA", "FTAL.TA", "FIBIH.TA", "KEN.TA", "MTRX.TA",
	"GILT.TA", "TRPZ.TA", "ISRA.TA", "DIMRI.TA", "ILCO.TA", "NOFR.TA", "ELTR.TA", "RATI.TA", "CEL.TA", "AURA.TA",
	"ISCD.TA", "ARYT.TA", "ASHG.TA", "MSKE.TA", "FORTY.TA", "MAXO.TA", "RMLI.TA", "ISCN.TA", "HLAN.TA", "ARPT.TA",
	"ORL.TA", "ONE.TA", "ISRO.TA", "INRM.TA", "NYAX.TA", "SLARL.TA", "PRTC.TA", "EQTL.TA", "FOX.TA", "GVYM.TA",
	"MISH.TA", "DANE.TA", "ARGO.TA", "BLSR.TA", "DNYA.TA", "SMT.TA", "IDIN.TA", "DELG.TA", "MRIN.TA", "ELCRE.TA",
	"NTML.TA", "AFRE.TA", "LAPD.TA", "AZRM.TA", "GCT.TA", "IBI.TA", "YHNF.TA", "CRSM.TA", "PRSK.TA", "VRDS.TA",
	"TMRP.TA", "ELCO.TA", "ACRO.TA", "ISRS.TA", "DLEA.TA", "RMON.TA", "AMRM.TA", "ECP.TA", "RTLS.TA", "ISHO.TA",
	"DUNI.TA", "OPK.TA", "PTBL.TA", "SBEN.TA", "ACKR.TA", "VILR.TA", "CRSR.TA", "ISHI.TA", "PLRM.TA", "DLTI.TA",
	"MLTM.TA", "AMPA.TA", "UNMI.TA",
}

var sp500Symbols = []string{
	"MMM", "AOS", "ABT", "ABBV", "ACN", "ADBE", "AMD", "AES", "AFL", "A",
	"APD", "ABNB", "AKAM", "ALB", "ARE", "ALGN", "ALLE", "LNT", "ALL", "GOOGL",
	"GOOG", "MO", "AMZN", "AMCR", "AEE", "AEP", "AXP", "AIG", "AMT", "AWK",
	"AMP", "AME", "AMGN", "APH", "ADI", "AON", "APA", "APO", "AAPL", "AMAT",
	"APP", "APTV", "ACGL", "ADM", "ARES", "ANET", "AJG", "AIZ", "T", "ATO",
	"ADSK", "ADP", "AZO", "AVB", "AVY", "AXON", "BKR", "BALL", "BAC", "BAX",
	"BDX", "BRK-B", "BBY", "TECH", "BIIB", "BLK", "BX", "XYZ", "BK", "BA",
	"BKNG", "BSX", "BMY", "AVGO", "BR", "BRO", "BF-B", "BLDR", "BG", "BXP",
	"CHRW", "CDNS", "CPT", "CPB", "COF", "CAH", "CCL", "CARR", "CVNA", "CASY",
	"CAT", "CBOE", "CBRE", "CDW", "COR", "CNC", "CNP", "CF", "CRL", "SCHW",
	"CHTR", "CVX", "CMG", "CB", "CHD", "CIEN", "CI", "CINF", "CTAS", "CSCO",
	"C", "CFG", "CLX", "CME", "CMS", "KO", "CTSH", "COHR", "COIN", "CL",
	"CMCSA", "FIX", "CAG", "COP", "ED", "STZ", "CEG", "COO", "CPRT", "GLW",
	"CPAY", "CTVA", "CSGP", "COST", "CTRA", "CRH", "CRWD", "CCI", "CSX", "CMI",
	"CVS", "DHR", "DRI", "DDOG", "DVA", "DECK", "DE", "DELL", "DAL", "DVN",
	"DXCM", "FANG", "DLR", "DG", "DLTR", "D", "DPZ", "DASH", "DOV", "DOW",
	"DHI", "DTE", "DUK", "DD", "ETN", "EBAY", "SATS", "ECL", "EIX", "EW",
	"EA", "ELV", "EME", "EMR", "ETR", "EOG", "EPAM", "EQT", "EFX", "EQIX",
	"EQR", "ERIE", "ESS", "EL", "EG", "EVRG", "ES", "EXC", "EXE", "EXPE",
	"EXPD", "EXR", "XOM", "FFIV", "FDS", "FICO", "FAST", "FRT", "FDX", "FIS",
	"FITB", "FSLR", "FE", "FISV", "F", "FTNT", "FTV", "FOXA", "FOX", "BEN",
	"FCX", "GRMN", "IT", "GE", "GEHC", "GEV", "GEN", "GNRC", "GD", "GIS",
	"GM", "GPC", "GILD", "GPN", "GL", "GDDY", "GS", "HAL", "HIG", "HAS",
	"HCA", "DOC", "HSIC", "HSY", "HPE", "HLT", "HD", "HON", "HRL", "HST",
	"HWM", "HPQ", "HUBB", "HUM", "HBAN", "HII", "IBM", "IEX", "IDXX", "ITW",
	"INCY", "IR", "PODD", "INTC", "IBKR", "ICE", "IFF", "IP", "INTU", "ISRG",
	"IVZ", "INVH", "IQV", "IRM", "JBHT", "JBL", "JKHY", "J", "JNJ", "JCI",
	"JPM", "KVUE", "KDP", "KEY", "KEYS", "KMB", "KIM", "KMI", "KKR", "KLAC",
	"KHC", "KR", "LHX", "LH", "LRCX", "LVS", "LDOS", "LEN", "LII", "LLY",
	"LIN", "LYV", "LMT", "L", "LOW", "LULU", "LITE", "LYB", "MTB", "MPC",
	"MAR", "MRSH", "MLM", "MAS", "MA", "MKC", "MCD", "MCK", "MDT", "MRK",
	"META", "MET", "MTD", "MGM", "MCHP", "MU", "MSFT", "MAA", "MRNA", "TAP",
	"MDLZ", "MPWR", "MNST", "MCO", "MS", "MOS", "MSI", "MSCI", "NDAQ", "NTAP",
	"NFLX", "NEM", "NWSA", "NWS", "NEE", "NKE", "NI", "NDSN", "NSC", "NTRS",
	"NOC", "NCLH", "NRG", "NUE", "NVDA", "NVR", "NXPI", "ORLY", "OXY", "ODFL",
	"OMC", "ON", "OKE", "ORCL", "OTIS", "PCAR", "PKG", "PLTR", "PANW", "PSKY",
	"PH", "PAYX", "PYPL", "PNR", "PEP", "PFE", "PCG", "PM", "PSX", "PNW",
	"PNC", "POOL", "PPG", "PPL", "PFG", "PG", "PGR", "PLD", "PRU", "PEG",
	"PTC", "PSA", "PHM", "PWR", "QCOM", "DGX", "Q", "RL", "RJF", "RTX",
	"O", "REG", "REGN", "RF", "RSG", "RMD", "RVTY", "HOOD", "ROK", "ROL",
	"ROP", "ROST", "RCL", "SPGI", "CRM", "SNDK", "SBAC", "SLB", "STX", "SRE",
	"NOW", "SHW", "SPG", "SWKS", "SJM", "SW", "SNA", "SOLV", "SO", "LUV",
	"SWK", "SBUX", "STT", "STLD", "STE", "SYK", "SMCI", "SYF", "SNPS", "SYY",
	"TMUS", "TROW", "TTWO", "TPR", "TRGP", "TGT", "TEL", "TDY", "TER", "TSLA",
	"TXN", "TPL", "TXT", "TMO", "TJX", "TKO", "TTD", "TSCO", "TT", "TDG",
	"TRV", "TRMB", "TFC", "TYL", "TSN", "USB", "UBER", "UDR", "ULTA", "UNP",
	"UAL", "UPS", "URI", "UNH", "UHS", "VLO", "VTR", "VLTO", "VRSN", "VRSK",
	"VZ", "VRTX", "VRT", "VTRS", "VICI", "V", "VST", "VMC", "WRB", "GWW",
	"WAB", "WMT", "DIS", "WBD", "WM", "WAT", "WEC", "WFC", "WELL", "WST",
	"WDC", "WY", "WSM", "WMB", "WTW", "WDAY", "WYNN", "XEL", "XYL", "YUM",
	"ZBRA", "ZBH", "ZTS",
}

var nasdaq100Symbols = []string{
	"ADBE", "AMD", "ABNB", "ALNY", "GOOGL", "GOOG", "AMZN", "AEP", "AMGN", "ADI",
	"AAPL", "AMAT", "APP", "ARM", "ASML", "TEAM", "ADSK", "ADP", "AXON", "BKR",
	"BKNG", "AVGO", "CDNS", "CHTR", "CTAS", "CSCO", "CCEP", "CTSH", "CMCSA", "CEG",
	"CPRT", "CSGP", "COST", "CRWD", "CSX", "DDOG", "DXCM", "FANG", "DASH", "EA",
	"EXC", "FAST", "FER", "FTNT", "GEHC", "GILD", "HON", "IDXX", "INSM", "INTC",
	"INTU", "ISRG", "KDP", "KLAC", "KHC", "LRCX", "LIN", "MAR", "MRVL", "MELI",
	"META", "MCHP", "MU", "MSFT", "MSTR", "MDLZ", "MPWR", "MNST", "NFLX", "NVDA",
	"NXPI", "ORLY", "ODFL", "PCAR", "PLTR", "PANW", "PAYX", "PYPL", "PDD", "PEP",
	"QCOM", "REGN", "ROP", "ROST", "STX", "SHOP", "SBUX", "SNPS", "TMUS", "TTWO",
	"TSLA", "TXN", "TRI", "VRSK", "VRTX", "WMT", "WBD", "WDC", "WDAY", "XEL",
	"ZS",
}

var legacyDefaultSymbols = []string{
	"AAPL", "MSFT", "GOOGL", "AMZN", "META", "NVDA", "TSLA", "AMD", "INTC", "CRM",
	"ORCL", "ADBE", "CSCO", "IBM", "QCOM", "TXN", "AVGO", "NOW", "SNOW", "PLTR",
	"JPM", "BAC", "WFC", "GS", "MS", "C", "BLK", "SCHW", "AXP", "V",
	"MA", "PYPL", "SQ", "COIN", "HOOD", "JNJ", "UNH", "PFE", "ABBV", "MRK",
	"LLY", "TMO", "ABT", "BMY", "AMGN", "GILD", "REGN", "VRTX", "MRNA", "BIIB",
	"WMT", "HD", "COST", "NKE", "MCD", "SBUX", "TGT", "LOW", "TJX", "LULU",
	"DIS", "NFLX", "CMCSA", "PEP", "KO", "CAT", "DE", "BA", "HON", "UPS",
	"FDX", "GE", "MMM", "LMT", "RTX", "XOM", "CVX", "COP", "SLB", "EOG",
	"MPC", "VLO", "PSX", "OXY", "DVN", "AMT", "PLD", "CCI", "EQIX", "SPG",
	"O", "WELL", "DLR", "AVB", "EQR", "NEE", "DUK", "SO", "D", "AEP",
	"EXC", "SRE", "XEL", "PEG", "ED", "T", "VZ", "TMUS", "CHTR", "LIN",
	"APD", "SHW", "ECL", "FCX", "NEM", "NUE",
}

// buildDefaultStockUniverse combines index lists with legacy coverage and deduplicates.
func buildDefaultStockUniverse() []string {
	combined := make([]string, 0, len(taIndexSymbols)+len(sp500Symbols)+len(nasdaq100Symbols)+len(legacyDefaultSymbols))
	combined = append(combined, taIndexSymbols...)
	combined = append(combined, sp500Symbols...)
	combined = append(combined, nasdaq100Symbols...)
	combined = append(combined, legacyDefaultSymbols...)

	seen := make(map[string]struct{}, len(combined))
	result := make([]string, 0, len(combined))
	for _, symbol := range combined {
		if _, ok := seen[symbol]; ok {
			continue
		}
		seen[symbol] = struct{}{}
		result = append(result, symbol)
	}
	return result
}
