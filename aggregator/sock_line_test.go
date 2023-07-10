package aggregator

import (
	"testing"
)

func TestSocketLine(t *testing.T) {
	sockLine := &SocketLine{
		Values: []TimestampedSocket{},
	}

	tsList := []uint64{

		33805065332163,
		33805065990716,
		33805066400606,
		33805066937463,
		33805067507004,
		33805068082621,
		33805068543449,
		33805069106660,
		33805069572630,
		33805070210774,
		33805070772370,
		33805071162619,
		33805071625600,
		33805073482028,
		33805073841739,
		33805074342888,
		33805074573808,
		33805075080976,
		33805075542978,
		33805076175534,
		33807002886899,
		33807004747817,
		33815077484050,
		33815078136583,
		33815078684015,
		33815079337782,
		33815079970725,
		33815080631457,
		33815081298490,
		33815081879427,
		33815082380024,
		33815084675593,
		33815086019234,
		33815086874036,
		33815087304858,
		33815088044606,
		33815088458212,
		33815089096509,
		33815089520695,
		33815090099952,
		33815090627646,
		33815091218320,
		33817002383633,
		33817004246338,
		33825091976123,
		33825092401903,
		33825092874397,
		33825093276700,
		33825093730185,
		33825094133024,
		33825094693337,
		33825095077622,
		33825095383712,
		33825095879835,
		33825096144660,
		33825097444170,
		33825097893511,
		33825098209615,
		33825098702565,
		33825099321611,
		33825102631064,
		33825103159797,
		33825103734780,
		33825104275145,
		33827002408011,
		33827004478008,
		33835105262438,
		33835105708889,
		33835106419631,
		33835106868158,
		33835108439083,
		33835108948978,
		33835109377337,
		33835109733899,
		33835110431996,
		33835110783593,
		33835111359964,
		33835111773541,
		33835112152140,
		33835112648775,
		33835113588946,
		33835114114962,
		33835116464267,
		33835116796833,
		33837002465570,
		33837005363575,
		33845117927928,
		33845118418834,
		33845118902184,
		33845119324753,
		33845119761634,
		33845120150932,
		33845120373773,
		33845120774746,
		33845120977225,
		33845121256138,
		33845121511775,
		33845121886526,
		33845122179563,
		33845122465899,
		33845122683308,
		33845122957336,
		33845123193483,
		33845123768007,
		33845124081620,
		33845124484199,
		33847002829796,
		33847004491036,
		33855126637018,
		33855127292811,
		33855127797136,
		33855129333727,
		33855129663006,
		33855130174887,
		33855130432133,
		33855131023596,
		33855131441985,
		33855131877348,
		33855132203368,
		33855132703031,
		33855133117599,
		33855133529811,
		33855133948340,
		33855134338365,
		33855135349858,
		33855136527049,
		33855137005564,
		33855138635975,
		33857002161562,
		33857003888105,
		33865139668514,
		33865140247948,
		33865140546565,
		33865140969319,
		33865141333634,
		33865141819179,
		33865142625886,
		33865143033130,
		33865143380604,
		33865143720225,
		33865143984814,
		33865144402045,
		33865144842268,
		33865146494654,
		33865146881021,
		33865147411956,
		33865147829404,
		33865148236623,
		33865148610611,
		33865149382097,
		33867002484049,
		33867004237247,
		33875151570756,
		33875152030350,
		33875152401364,
		33875152766760,
		33875153000748,
		33875153385347,
		33875153887830,
		33875154493200,
		33875154647639,
		33875155158138,
		33875155394853,
		33875155789707,
		33875156116705,
		33875157056085,
		33875157584857,
		33875157967095,
		33875160067274,
		33875160477500,
		33875162322749,
		33875162716805,
		33877002080006,
		33877004220942,
		33885163902686,
		33885164391873,
		33885164770207,
		33885165199364,
		33885165506387,
		33885167132542,
		33885167888718,
		33885168287316,
		33885168788429,
		33885169316003,
		33885169837989,
		33885170232629,
		33885172006997,
		33885172532136,
		33885173057385,
		33885173432203,
		33885173672241,
		33885174261255,
		33885174516480,
		33885174949773,
		33887002863735,
		33887004543465,
		33905182247187,
		33905182681540,
		33905182921150,
		33905183331246,
		33905183684506,
		33905184041291,
		33905184731352,
		33905185075193,
		33905185226477,
		33905185666810,
		33905186155628,
		33905186467885,
		33905187321782,
		33905187613986,
		33905188216655,
		33905188569287,
		33905188845524,
		33905192931995,
		33905193336688,
		33907002789359,
		33907004763240,
		33915195311234,
		33915195884981,
		33915196115400,
		33915196521505,
		33915196709821,
		33915197080933,
		33915197288006,
		33915200972646,
		33915201344393,
		33915201757715,
		33915201958333,
		33915202331043,
		33915202552713,
		33915202904426,
		33915203125022,
		33915203473423,
		33915203684961,
		33915204020392,
		33915204279152,
		33915204628464,
		33916758058783,
		33916758549119,
		33916758808545,
		33916759265974,
		33916759655277,
		33916760390180,
		33917002238669,
		33917004229882,
		33925205488211,
		33925205915908,
		33925206174793,
		33925206473701,
		33925206656234,
		33925207065308,
		33925207301307,
		33925207820562,
		33925208000464,
		33925208288320,
		33925208629623,
		33925208909228,
		33925209145687,
		33925209411134,
		33925209611028,
		33925209860305,
		33925210033376,
		33925218073839,
		33925218412591,
		33925218725390,
		33927003077356,
		33927005428013,
		33935220106305,
		33935220589599,
		33935220836026,
		33935221143996,
		33935221367883,
		33935221657852,
		33935222018048,
		33935222361191,
		33935222576305,
		33935222986404,
		33935223340999,
		33935223656483,
		33935223864611,
		33935224179065,
		33935224428683,
		33935224731612,
		33935224940794,
		33935225227144,
		33935225412521,
		33935225677159,
		33937002928428,
		33937005112199,
		33945226651486,
		33945227108602,
		33945227458231,
		33945227942270,
		33945228262614,
		33945228767843,
		33945229191007,
		33945229775947,
		33945229987129,
		33945230576079,
		33945230797176,
		33945231235604,
		33945232859491,
		33945234961387,
		33945235683085,
		33945236269094,
		33945236611501,
		33947002331172,
		33947004517045,
	}

	for _, ts := range tsList {
		sockLine.AddValue(ts, &SockInfo{
			EstablishTime: ts,
		})
	}

	// 21615006690038,
	// 21615006748038
	// 21615006991934,

	// 3383510 6868158,
	// 3383510 7729129
	// 3383510 8439083,

	si, err := sockLine.GetValue(33835107729129)
	if err != nil || si == nil {
		t.Fatalf("unexpected error: %v", err)
		t.Fail()
		return
	}

}
