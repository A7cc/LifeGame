export namespace core {
	
	export class Announce {
	    announceins: string[];
	    announceout: string[];
	    announcecompany: string[];
	    announcegame: string[];
	    announcehealthy: string[];
	
	    static createFrom(source: any = {}) {
	        return new Announce(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.announceins = source["announceins"];
	        this.announceout = source["announceout"];
	        this.announcecompany = source["announcecompany"];
	        this.announcegame = source["announcegame"];
	        this.announcehealthy = source["announcehealthy"];
	    }
	}
	export class AntiqueInfo {
	    aiid: number;
	    ainame: string;
	    aiprice: number;
	    aiidisplay: number;
	    aiamaterial: number;
	    aiacondition: number;
	    aiprice_max: number;
	    aiatime: number;
	    aiimg: string;
	    aidesc: string;
	    ailevel: number;
	
	    static createFrom(source: any = {}) {
	        return new AntiqueInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.aiid = source["aiid"];
	        this.ainame = source["ainame"];
	        this.aiprice = source["aiprice"];
	        this.aiidisplay = source["aiidisplay"];
	        this.aiamaterial = source["aiamaterial"];
	        this.aiacondition = source["aiacondition"];
	        this.aiprice_max = source["aiprice_max"];
	        this.aiatime = source["aiatime"];
	        this.aiimg = source["aiimg"];
	        this.aidesc = source["aidesc"];
	        this.ailevel = source["ailevel"];
	    }
	}
	export class BankTask {
	    taskid: number;
	    taskname: string;
	    taskdesc: string;
	    tasktype: string;
	    targetvalue: number;
	    reward: number;
	
	    static createFrom(source: any = {}) {
	        return new BankTask(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.taskid = source["taskid"];
	        this.taskname = source["taskname"];
	        this.taskdesc = source["taskdesc"];
	        this.tasktype = source["tasktype"];
	        this.targetvalue = source["targetvalue"];
	        this.reward = source["reward"];
	    }
	}
	export class BankTaskStats {
	    netbankflow: number;
	    depositamount: number;
	    depositcount: number;
	    withdrawamount: number;
	    withdrawcount: number;
	    loanamount: number;
	    workcount: number;
	    completedtasks: number[];
	    claimedtasks: number[];
	    currenttasks: BankTask[];
	
	    static createFrom(source: any = {}) {
	        return new BankTaskStats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.netbankflow = source["netbankflow"];
	        this.depositamount = source["depositamount"];
	        this.depositcount = source["depositcount"];
	        this.withdrawamount = source["withdrawamount"];
	        this.withdrawcount = source["withdrawcount"];
	        this.loanamount = source["loanamount"];
	        this.workcount = source["workcount"];
	        this.completedtasks = source["completedtasks"];
	        this.claimedtasks = source["claimedtasks"];
	        this.currenttasks = this.convertValues(source["currenttasks"], BankTask);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class CarInfo {
	    ciid: number;
	    ciname: string;
	    ciprice: number;
	    cihealth: number;
	    cifame: number;
	    ciimg: string;
	
	    static createFrom(source: any = {}) {
	        return new CarInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ciid = source["ciid"];
	        this.ciname = source["ciname"];
	        this.ciprice = source["ciprice"];
	        this.cihealth = source["cihealth"];
	        this.cifame = source["cifame"];
	        this.ciimg = source["ciimg"];
	    }
	}
	export class CompanyInfo {
	    ciname: string;
	    ciprice: number;
	    cirisk: number;
	    ciprofit: number;
	    citime: number;
	    cistatus: boolean;
	
	    static createFrom(source: any = {}) {
	        return new CompanyInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ciname = source["ciname"];
	        this.ciprice = source["ciprice"];
	        this.cirisk = source["cirisk"];
	        this.ciprofit = source["ciprofit"];
	        this.citime = source["citime"];
	        this.cistatus = source["cistatus"];
	    }
	}
	export class MeetCondition {
	    ctype: string;
	    cvalue: number;
	    ctarget: string;
	
	    static createFrom(source: any = {}) {
	        return new MeetCondition(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ctype = source["ctype"];
	        this.cvalue = source["cvalue"];
	        this.ctarget = source["ctarget"];
	    }
	}
	export class DatingInfo {
	    did: number;
	    dname: string;
	    dimage: string;
	    dage: number;
	    dsex: boolean;
	    dnationality: string;
	    doccup: string;
	    ddesc: string;
	    dcost: number;
	    dmeetconditions: MeetCondition[];
	    dgifts: string[];
	    dlocations: string[];
	    dmeetscene: string;
	    dunlocked: boolean;
	    daffinitylevel: string;
	
	    static createFrom(source: any = {}) {
	        return new DatingInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.did = source["did"];
	        this.dname = source["dname"];
	        this.dimage = source["dimage"];
	        this.dage = source["dage"];
	        this.dsex = source["dsex"];
	        this.dnationality = source["dnationality"];
	        this.doccup = source["doccup"];
	        this.ddesc = source["ddesc"];
	        this.dcost = source["dcost"];
	        this.dmeetconditions = this.convertValues(source["dmeetconditions"], MeetCondition);
	        this.dgifts = source["dgifts"];
	        this.dlocations = source["dlocations"];
	        this.dmeetscene = source["dmeetscene"];
	        this.dunlocked = source["dunlocked"];
	        this.daffinitylevel = source["daffinitylevel"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class DifficultyConfig {
	    level: number;
	    name: string;
	    initmoney: number;
	    healthbonus: number;
	    bankruptRate: number;
	    antiquefake: number;
	    description: string;
	
	    static createFrom(source: any = {}) {
	        return new DifficultyConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.level = source["level"];
	        this.name = source["name"];
	        this.initmoney = source["initmoney"];
	        this.healthbonus = source["healthbonus"];
	        this.bankruptRate = source["bankruptRate"];
	        this.antiquefake = source["antiquefake"];
	        this.description = source["description"];
	    }
	}
	export class HouseInfo {
	    hiid: number;
	    hiname: string;
	    hiprice: number;
	    hihealth: number;
	    hifame: number;
	    hiimg: string;
	
	    static createFrom(source: any = {}) {
	        return new HouseInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.hiid = source["hiid"];
	        this.hiname = source["hiname"];
	        this.hiprice = source["hiprice"];
	        this.hihealth = source["hihealth"];
	        this.hifame = source["hifame"];
	        this.hiimg = source["hiimg"];
	    }
	}
	export class StockInfo {
	    siid: number;
	    siname: string;
	    siprice: number;
	    sirisk: number;
	    sihistory: number[];
	    siklinehistory: number[][];
	    sistatus: string;
	
	    static createFrom(source: any = {}) {
	        return new StockInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.siid = source["siid"];
	        this.siname = source["siname"];
	        this.siprice = source["siprice"];
	        this.sirisk = source["sirisk"];
	        this.sihistory = source["sihistory"];
	        this.siklinehistory = source["siklinehistory"];
	        this.sistatus = source["sistatus"];
	    }
	}
	export class MaxHoldNum {
	    mdholdnum: number;
	    mfholdnum: number;
	    mcholdnum: number;
	    maholdnum: number;
	    mfaholdnum: number;
	    miholdnum: number;
	    mwroundnum: number;
	    mgroundnum: number;
	    mmroundnum: number;
	    msroundnum: number;
	    maroundnum: number;
	
	    static createFrom(source: any = {}) {
	        return new MaxHoldNum(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.mdholdnum = source["mdholdnum"];
	        this.mfholdnum = source["mfholdnum"];
	        this.mcholdnum = source["mcholdnum"];
	        this.maholdnum = source["maholdnum"];
	        this.mfaholdnum = source["mfaholdnum"];
	        this.miholdnum = source["miholdnum"];
	        this.mwroundnum = source["mwroundnum"];
	        this.mgroundnum = source["mgroundnum"];
	        this.mmroundnum = source["mmroundnum"];
	        this.msroundnum = source["msroundnum"];
	        this.maroundnum = source["maroundnum"];
	    }
	}
	export class ItemInfo {
	    iiname: string;
	    iiprice: number;
	    iieffect: string;
	    iidisplay: boolean;
	
	    static createFrom(source: any = {}) {
	        return new ItemInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.iiname = source["iiname"];
	        this.iiprice = source["iiprice"];
	        this.iieffect = source["iieffect"];
	        this.iidisplay = source["iidisplay"];
	    }
	}
	export class Game {
	    gid: number;
	    gname: string;
	    gtime: number;
	    gdifficulty: number;
	    giteminsinfo: Record<number, ItemInfo>;
	    gitemoutinfo: Record<number, ItemInfo>;
	    gcompanyinfo: Record<number, CompanyInfo>;
	    gantiqueinfo: AntiqueInfo;
	    gmaxholdnum: MaxHoldNum;
	    gstockinfo: StockInfo[];
	    gstocknews: string[];
	    gstockupdatecount: number;
	    gbanktaskstats: BankTaskStats;
	    gdatinginfo: DatingInfo[];
	    ghouseinfo: Record<number, HouseInfo>;
	    gcarinfo: Record<number, CarInfo>;
	
	    static createFrom(source: any = {}) {
	        return new Game(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.gid = source["gid"];
	        this.gname = source["gname"];
	        this.gtime = source["gtime"];
	        this.gdifficulty = source["gdifficulty"];
	        this.giteminsinfo = this.convertValues(source["giteminsinfo"], ItemInfo, true);
	        this.gitemoutinfo = this.convertValues(source["gitemoutinfo"], ItemInfo, true);
	        this.gcompanyinfo = this.convertValues(source["gcompanyinfo"], CompanyInfo, true);
	        this.gantiqueinfo = this.convertValues(source["gantiqueinfo"], AntiqueInfo);
	        this.gmaxholdnum = this.convertValues(source["gmaxholdnum"], MaxHoldNum);
	        this.gstockinfo = this.convertValues(source["gstockinfo"], StockInfo);
	        this.gstocknews = source["gstocknews"];
	        this.gstockupdatecount = source["gstockupdatecount"];
	        this.gbanktaskstats = this.convertValues(source["gbanktaskstats"], BankTaskStats);
	        this.gdatinginfo = this.convertValues(source["gdatinginfo"], DatingInfo);
	        this.ghouseinfo = this.convertValues(source["ghouseinfo"], HouseInfo, true);
	        this.gcarinfo = this.convertValues(source["gcarinfo"], CarInfo, true);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	
	
	
	export class MiniGameRecord {
	    mgrtype: string;
	    playcount: number;
	    wincount: number;
	
	    static createFrom(source: any = {}) {
	        return new MiniGameRecord(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.mgrtype = source["mgrtype"];
	        this.playcount = source["playcount"];
	        this.wincount = source["wincount"];
	    }
	}
	export class Opportunity {
	    ownum: number;
	    ognum: number;
	    omnum: number;
	    osnum: number;
	    oanum: number;
	
	    static createFrom(source: any = {}) {
	        return new Opportunity(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ownum = source["ownum"];
	        this.ognum = source["ognum"];
	        this.omnum = source["omnum"];
	        this.osnum = source["osnum"];
	        this.oanum = source["oanum"];
	    }
	}
	
	export class UCompanyInfo {
	    ucompanyname: string;
	    ucompanyholdtime: number;
	    ucompanycostprice: number;
	    ucompanynum: number;
	
	    static createFrom(source: any = {}) {
	        return new UCompanyInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ucompanyname = source["ucompanyname"];
	        this.ucompanyholdtime = source["ucompanyholdtime"];
	        this.ucompanycostprice = source["ucompanycostprice"];
	        this.ucompanynum = source["ucompanynum"];
	    }
	}
	export class UDiseaseInfo {
	    udname: string;
	    udtype: string;
	    usymptoms: string;
	    uhealthimpact: number;
	    uupgradedays: number;
	    utreatments: string[];
	    udseverity: number;
	    udtime: number;
	
	    static createFrom(source: any = {}) {
	        return new UDiseaseInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.udname = source["udname"];
	        this.udtype = source["udtype"];
	        this.usymptoms = source["usymptoms"];
	        this.uhealthimpact = source["uhealthimpact"];
	        this.uupgradedays = source["uupgradedays"];
	        this.utreatments = source["utreatments"];
	        this.udseverity = source["udseverity"];
	        this.udtime = source["udtime"];
	    }
	}
	export class UItemInfo {
	    uicostprice: number;
	    uitemnum: number;
	
	    static createFrom(source: any = {}) {
	        return new UItemInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.uicostprice = source["uicostprice"];
	        this.uitemnum = source["uitemnum"];
	    }
	}
	export class UserDatingInfo {
	    ddatingid: number;
	    dname: string;
	    daffinity: number;
	    dcount: number;
	    dgiftcount: number;
	    dstatus: string;
	
	    static createFrom(source: any = {}) {
	        return new UserDatingInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ddatingid = source["ddatingid"];
	        this.dname = source["dname"];
	        this.daffinity = source["daffinity"];
	        this.dcount = source["dcount"];
	        this.dgiftcount = source["dgiftcount"];
	        this.dstatus = source["dstatus"];
	    }
	}
	export class UserStockInfo {
	    usname: string;
	    usprice_init: number;
	    usnum: number;
	    usprofit: number;
	
	    static createFrom(source: any = {}) {
	        return new UserStockInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.usname = source["usname"];
	        this.usprice_init = source["usprice_init"];
	        this.usnum = source["usnum"];
	        this.usprofit = source["usprofit"];
	    }
	}
	export class User {
	    uid: number;
	    uname: string;
	    usex: boolean;
	    uage: number;
	    uimmunity: number;
	    ucriticalhealthyears: number;
	    udiseases: Record<number, UDiseaseInfo>;
	    ufame: number;
	    uopportunity: Opportunity;
	    uassets: number;
	    ucash: number;
	    ubank: number;
	    uloan: number;
	    uloanoverdue: number;
	    uitemins: Record<number, UItemInfo>;
	    uitemout: Record<number, UItemInfo>;
	    uantique: AntiqueInfo[];
	    ucompany: Record<number, UCompanyInfo>;
	    ustock: Record<number, UserStockInfo>;
	    ustockprofit: number;
	    udating: Record<number, UserDatingInfo>;
	    umarrieddatingid: number;
	    ucar: Record<number, boolean>;
	    uhouse: Record<number, boolean>;
	    uminigamerecords: Record<string, MiniGameRecord>;
	
	    static createFrom(source: any = {}) {
	        return new User(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.uid = source["uid"];
	        this.uname = source["uname"];
	        this.usex = source["usex"];
	        this.uage = source["uage"];
	        this.uimmunity = source["uimmunity"];
	        this.ucriticalhealthyears = source["ucriticalhealthyears"];
	        this.udiseases = this.convertValues(source["udiseases"], UDiseaseInfo, true);
	        this.ufame = source["ufame"];
	        this.uopportunity = this.convertValues(source["uopportunity"], Opportunity);
	        this.uassets = source["uassets"];
	        this.ucash = source["ucash"];
	        this.ubank = source["ubank"];
	        this.uloan = source["uloan"];
	        this.uloanoverdue = source["uloanoverdue"];
	        this.uitemins = this.convertValues(source["uitemins"], UItemInfo, true);
	        this.uitemout = this.convertValues(source["uitemout"], UItemInfo, true);
	        this.uantique = this.convertValues(source["uantique"], AntiqueInfo);
	        this.ucompany = this.convertValues(source["ucompany"], UCompanyInfo, true);
	        this.ustock = this.convertValues(source["ustock"], UserStockInfo, true);
	        this.ustockprofit = source["ustockprofit"];
	        this.udating = this.convertValues(source["udating"], UserDatingInfo, true);
	        this.umarrieddatingid = source["umarrieddatingid"];
	        this.ucar = source["ucar"];
	        this.uhouse = source["uhouse"];
	        this.uminigamerecords = this.convertValues(source["uminigamerecords"], MiniGameRecord, true);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	

}

export namespace services {
	
	export class AppStartupStatus {
	    ready: boolean;
	    stage: string;
	    error: string;
	    updatedAt: number;
	
	    static createFrom(source: any = {}) {
	        return new AppStartupStatus(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ready = source["ready"];
	        this.stage = source["stage"];
	        this.error = source["error"];
	        this.updatedAt = source["updatedAt"];
	    }
	}
	export class BasicResponse {
	    code: number;
	    msg: string;
	
	    static createFrom(source: any = {}) {
	        return new BasicResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.code = source["code"];
	        this.msg = source["msg"];
	    }
	}
	export class DatingGiftOption {
	    name: string;
	    cost: number;
	    preferredeffect: number;
	    riskyeffect: number;
	
	    static createFrom(source: any = {}) {
	        return new DatingGiftOption(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.cost = source["cost"];
	        this.preferredeffect = source["preferredeffect"];
	        this.riskyeffect = source["riskyeffect"];
	    }
	}
	export class DatingGiftOptionsResponse {
	    code: number;
	    msg: string;
	    options?: DatingGiftOption[];
	
	    static createFrom(source: any = {}) {
	        return new DatingGiftOptionsResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.code = source["code"];
	        this.msg = source["msg"];
	        this.options = this.convertValues(source["options"], DatingGiftOption);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class DatingInteractionResponse {
	    code: number;
	    msg: string;
	    interaction?: string;
	    outcome?: string;
	    label?: string;
	    event?: string;
	    outfit?: string;
	    outfitvariant?: string;
	    outfitimage?: string;
	    affinitychange: number;
	    datinginfo?: core.UserDatingInfo;
	    userinfo?: core.User;
	
	    static createFrom(source: any = {}) {
	        return new DatingInteractionResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.code = source["code"];
	        this.msg = source["msg"];
	        this.interaction = source["interaction"];
	        this.outcome = source["outcome"];
	        this.label = source["label"];
	        this.event = source["event"];
	        this.outfit = source["outfit"];
	        this.outfitvariant = source["outfitvariant"];
	        this.outfitimage = source["outfitimage"];
	        this.affinitychange = source["affinitychange"];
	        this.datinginfo = this.convertValues(source["datinginfo"], core.UserDatingInfo);
	        this.userinfo = this.convertValues(source["userinfo"], core.User);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class DatingListResponse {
	    code: number;
	    msg: string;
	    datinglist: core.DatingInfo[];
	    meetingscenes: string[];
	    userinfo?: core.User;
	
	    static createFrom(source: any = {}) {
	        return new DatingListResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.code = source["code"];
	        this.msg = source["msg"];
	        this.datinglist = this.convertValues(source["datinglist"], core.DatingInfo);
	        this.meetingscenes = source["meetingscenes"];
	        this.userinfo = this.convertValues(source["userinfo"], core.User);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class DatingRelationshipResponse {
	    code: number;
	    msg: string;
	    event?: string;
	    gift?: string;
	    giftcost?: number;
	    affinitychange?: number;
	    outcome?: string;
	    preferred: boolean;
	    success: boolean;
	    successrate?: number;
	    datinginfo?: core.UserDatingInfo;
	    userinfo?: core.User;
	
	    static createFrom(source: any = {}) {
	        return new DatingRelationshipResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.code = source["code"];
	        this.msg = source["msg"];
	        this.event = source["event"];
	        this.gift = source["gift"];
	        this.giftcost = source["giftcost"];
	        this.affinitychange = source["affinitychange"];
	        this.outcome = source["outcome"];
	        this.preferred = source["preferred"];
	        this.success = source["success"];
	        this.successrate = source["successrate"];
	        this.datinginfo = this.convertValues(source["datinginfo"], core.UserDatingInfo);
	        this.userinfo = this.convertValues(source["userinfo"], core.User);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class DatingSceneOutcome {
	    location: string;
	    category: string;
	    label: string;
	    event: string;
	    effect: string;
	    preferred: boolean;
	    preferenceratechange: number;
	    successrate: number;
	    rewardtier: string;
	    rewardmultiplier: number;
	    moment?: string;
	    momentlabel?: string;
	    momentevent?: string;
	    momentaffinitychange?: number;
	    famechange: number;
	    healthchange: number;
	    affinitychange: number;
	
	    static createFrom(source: any = {}) {
	        return new DatingSceneOutcome(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.location = source["location"];
	        this.category = source["category"];
	        this.label = source["label"];
	        this.event = source["event"];
	        this.effect = source["effect"];
	        this.preferred = source["preferred"];
	        this.preferenceratechange = source["preferenceratechange"];
	        this.successrate = source["successrate"];
	        this.rewardtier = source["rewardtier"];
	        this.rewardmultiplier = source["rewardmultiplier"];
	        this.moment = source["moment"];
	        this.momentlabel = source["momentlabel"];
	        this.momentevent = source["momentevent"];
	        this.momentaffinitychange = source["momentaffinitychange"];
	        this.famechange = source["famechange"];
	        this.healthchange = source["healthchange"];
	        this.affinitychange = source["affinitychange"];
	    }
	}
	export class DatingResultResponse {
	    code: number;
	    msg: string;
	    success: boolean;
	    famechange: number;
	    healthchange: number;
	    affinitychange: number;
	    scene?: DatingSceneOutcome;
	    datinginfo?: core.UserDatingInfo;
	    userinfo?: core.User;
	
	    static createFrom(source: any = {}) {
	        return new DatingResultResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.code = source["code"];
	        this.msg = source["msg"];
	        this.success = source["success"];
	        this.famechange = source["famechange"];
	        this.healthchange = source["healthchange"];
	        this.affinitychange = source["affinitychange"];
	        this.scene = this.convertValues(source["scene"], DatingSceneOutcome);
	        this.datinginfo = this.convertValues(source["datinginfo"], core.UserDatingInfo);
	        this.userinfo = this.convertValues(source["userinfo"], core.User);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class DatingSceneResponse {
	    code: number;
	    msg: string;
	    scene?: string;
	    met: core.DatingInfo[];
	    userinfo?: core.User;
	
	    static createFrom(source: any = {}) {
	        return new DatingSceneResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.code = source["code"];
	        this.msg = source["msg"];
	        this.scene = source["scene"];
	        this.met = this.convertValues(source["met"], core.DatingInfo);
	        this.userinfo = this.convertValues(source["userinfo"], core.User);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class GameEvaluation {
	    title: string;
	    description: string;
	    score: number;
	    wealthscore: number;
	    healthscore: number;
	    famescore: number;
	    agescore: number;
	    careerscore: number;
	    relationshipscore: number;
	    collectionscore: number;
	
	    static createFrom(source: any = {}) {
	        return new GameEvaluation(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.title = source["title"];
	        this.description = source["description"];
	        this.score = source["score"];
	        this.wealthscore = source["wealthscore"];
	        this.healthscore = source["healthscore"];
	        this.famescore = source["famescore"];
	        this.agescore = source["agescore"];
	        this.careerscore = source["careerscore"];
	        this.relationshipscore = source["relationshipscore"];
	        this.collectionscore = source["collectionscore"];
	    }
	}
	export class EvaluationResponse {
	    code: number;
	    msg: string;
	    evaluation?: GameEvaluation;
	
	    static createFrom(source: any = {}) {
	        return new EvaluationResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.code = source["code"];
	        this.msg = source["msg"];
	        this.evaluation = this.convertValues(source["evaluation"], GameEvaluation);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class GameStateResponse {
	    code: number;
	    msg: string;
	    gameinfo?: core.Game;
	    userinfo?: core.User;
	    announce?: core.Announce;
	    difficulty?: core.DifficultyConfig;
	    stockepoch?: string;
	    stockversion?: number;
	
	    static createFrom(source: any = {}) {
	        return new GameStateResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.code = source["code"];
	        this.msg = source["msg"];
	        this.gameinfo = this.convertValues(source["gameinfo"], core.Game);
	        this.userinfo = this.convertValues(source["userinfo"], core.User);
	        this.announce = this.convertValues(source["announce"], core.Announce);
	        this.difficulty = this.convertValues(source["difficulty"], core.DifficultyConfig);
	        this.stockepoch = source["stockepoch"];
	        this.stockversion = source["stockversion"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class SaveInfo {
	    id: number;
	    name: string;
	    created_at: string;
	    game_year: number;
	    save_version: number;
	
	    static createFrom(source: any = {}) {
	        return new SaveInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.created_at = source["created_at"];
	        this.game_year = source["game_year"];
	        this.save_version = source["save_version"];
	    }
	}
	export class ListSavesResponse {
	    code: number;
	    msg: string;
	    saves: SaveInfo[];
	
	    static createFrom(source: any = {}) {
	        return new ListSavesResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.code = source["code"];
	        this.msg = source["msg"];
	        this.saves = this.convertValues(source["saves"], SaveInfo);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class LoadGameResponse {
	    code: number;
	    msg: string;
	    gameinfo?: core.Game;
	    userinfo?: core.User;
	    announce?: core.Announce;
	    difficulty?: core.DifficultyConfig;
	    stockepoch?: string;
	    stockversion?: number;
	    saveversion?: number;
	
	    static createFrom(source: any = {}) {
	        return new LoadGameResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.code = source["code"];
	        this.msg = source["msg"];
	        this.gameinfo = this.convertValues(source["gameinfo"], core.Game);
	        this.userinfo = this.convertValues(source["userinfo"], core.User);
	        this.announce = this.convertValues(source["announce"], core.Announce);
	        this.difficulty = this.convertValues(source["difficulty"], core.DifficultyConfig);
	        this.stockepoch = source["stockepoch"];
	        this.stockversion = source["stockversion"];
	        this.saveversion = source["saveversion"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class SaveGameResponse {
	    code: number;
	    msg: string;
	    saveId?: number;
	    saveversion?: number;
	
	    static createFrom(source: any = {}) {
	        return new SaveGameResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.code = source["code"];
	        this.msg = source["msg"];
	        this.saveId = source["saveId"];
	        this.saveversion = source["saveversion"];
	    }
	}
	
	export class SpouseInteractionResponse {
	    code: number;
	    msg: string;
	    interaction?: string;
	    cost?: number;
	    affinitychange?: number;
	    healthchange?: number;
	    datinginfo?: core.UserDatingInfo;
	    userinfo?: core.User;
	
	    static createFrom(source: any = {}) {
	        return new SpouseInteractionResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.code = source["code"];
	        this.msg = source["msg"];
	        this.interaction = source["interaction"];
	        this.cost = source["cost"];
	        this.affinitychange = source["affinitychange"];
	        this.healthchange = source["healthchange"];
	        this.datinginfo = this.convertValues(source["datinginfo"], core.UserDatingInfo);
	        this.userinfo = this.convertValues(source["userinfo"], core.User);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class StartupResponse {
	    code: number;
	    msg: string;
	    status: AppStartupStatus;
	
	    static createFrom(source: any = {}) {
	        return new StartupResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.code = source["code"];
	        this.msg = source["msg"];
	        this.status = this.convertValues(source["status"], AppStartupStatus);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class StockUpdateResponse {
	    code: number;
	    msg: string;
	    stocks: core.StockInfo[];
	    news: string[];
	    epoch: string;
	    version: number;
	    updatedAt: number;
	    remaining: number;
	    marketclosed: boolean;
	
	    static createFrom(source: any = {}) {
	        return new StockUpdateResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.code = source["code"];
	        this.msg = source["msg"];
	        this.stocks = this.convertValues(source["stocks"], core.StockInfo);
	        this.news = source["news"];
	        this.epoch = source["epoch"];
	        this.version = source["version"];
	        this.updatedAt = source["updatedAt"];
	        this.remaining = source["remaining"];
	        this.marketclosed = source["marketclosed"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

