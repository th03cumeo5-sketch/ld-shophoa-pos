export namespace models {
	
	export class BanleTam {
	    khoa: number[];
	    lct: string;
	    // Go type: time
	    ngay: any;
	    sct: string;
	    loc: string;
	    makho: string;
	    quay: string;
	    ca: string;
	    thungan: string;
	    madt: string;
	    mavip: string;
	    soluong: number;
	    sotien: number;
	    chietkhau_mh: number;
	    tl_ck: number;
	    chietkhau: number;
	    thanhtien_chua_vat: number;
	    tien_thue: number;
	    tien_thue_cl: number;
	    thanhtoan: number;
	    pttt: string;
	    ttoan1: number;
	    ttoan2: number;
	    chuathoi: number;
	    printed: boolean;
	    print_no: number;
	    dgiai: string;
	    create_by: number;
	    // Go type: time
	    create_date: any;
	    update_by: number;
	    // Go type: time
	    update_date: any;
	
	    static createFrom(source: any = {}) {
	        return new BanleTam(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.khoa = source["khoa"];
	        this.lct = source["lct"];
	        this.ngay = this.convertValues(source["ngay"], null);
	        this.sct = source["sct"];
	        this.loc = source["loc"];
	        this.makho = source["makho"];
	        this.quay = source["quay"];
	        this.ca = source["ca"];
	        this.thungan = source["thungan"];
	        this.madt = source["madt"];
	        this.mavip = source["mavip"];
	        this.soluong = source["soluong"];
	        this.sotien = source["sotien"];
	        this.chietkhau_mh = source["chietkhau_mh"];
	        this.tl_ck = source["tl_ck"];
	        this.chietkhau = source["chietkhau"];
	        this.thanhtien_chua_vat = source["thanhtien_chua_vat"];
	        this.tien_thue = source["tien_thue"];
	        this.tien_thue_cl = source["tien_thue_cl"];
	        this.thanhtoan = source["thanhtoan"];
	        this.pttt = source["pttt"];
	        this.ttoan1 = source["ttoan1"];
	        this.ttoan2 = source["ttoan2"];
	        this.chuathoi = source["chuathoi"];
	        this.printed = source["printed"];
	        this.print_no = source["print_no"];
	        this.dgiai = source["dgiai"];
	        this.create_by = source["create_by"];
	        this.create_date = this.convertValues(source["create_date"], null);
	        this.update_by = source["update_by"];
	        this.update_date = this.convertValues(source["update_date"], null);
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
	export class BanleTamCT {
	    khoact: number[];
	    khoa: number[];
	    stt: number;
	    mavt: string;
	    tenvt: string;
	    dvt: string;
	    soluong: number;
	    dongia: number;
	    dongia_sauck: number;
	    sotien: number;
	    tl_ck: number;
	    tl_ck2: number;
	    tl_ck3: number;
	    tl_ck4: number;
	    tl_ck5: number;
	    chietkhau: number;
	    loaithue: string;
	    thue_suat: number;
	    tien_thue: number;
	    tien_thue_cl: number;
	    thanhtien_chua_vat: number;
	    thanhtoan: number;
	    ghichu: string;
	
	    static createFrom(source: any = {}) {
	        return new BanleTamCT(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.khoact = source["khoact"];
	        this.khoa = source["khoa"];
	        this.stt = source["stt"];
	        this.mavt = source["mavt"];
	        this.tenvt = source["tenvt"];
	        this.dvt = source["dvt"];
	        this.soluong = source["soluong"];
	        this.dongia = source["dongia"];
	        this.dongia_sauck = source["dongia_sauck"];
	        this.sotien = source["sotien"];
	        this.tl_ck = source["tl_ck"];
	        this.tl_ck2 = source["tl_ck2"];
	        this.tl_ck3 = source["tl_ck3"];
	        this.tl_ck4 = source["tl_ck4"];
	        this.tl_ck5 = source["tl_ck5"];
	        this.chietkhau = source["chietkhau"];
	        this.loaithue = source["loaithue"];
	        this.thue_suat = source["thue_suat"];
	        this.tien_thue = source["tien_thue"];
	        this.tien_thue_cl = source["tien_thue_cl"];
	        this.thanhtien_chua_vat = source["thanhtien_chua_vat"];
	        this.thanhtoan = source["thanhtoan"];
	        this.ghichu = source["ghichu"];
	    }
	}
	export class Category {
	    manhom: string;
	    tennhom: string;
	    count: number;
	
	    static createFrom(source: any = {}) {
	        return new Category(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.manhom = source["manhom"];
	        this.tennhom = source["tennhom"];
	        this.count = source["count"];
	    }
	}
	export class Customer {
	    madt: string;
	    tendt: string;
	    dienthoai: string;
	    diachi: string;
	    mst: string;
	    loai_vip: string;
	    sodu: number;
	
	    static createFrom(source: any = {}) {
	        return new Customer(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.madt = source["madt"];
	        this.tendt = source["tendt"];
	        this.dienthoai = source["dienthoai"];
	        this.diachi = source["diachi"];
	        this.mst = source["mst"];
	        this.loai_vip = source["loai_vip"];
	        this.sodu = source["sodu"];
	    }
	}
	export class ProductSearchResult {
	    mavt: string;
	    tenvt: string;
	    dvt: string;
	    giaban: number;
	    loaithue: string;
	    vat_ra: number;
	    tonkho: number;
	
	    static createFrom(source: any = {}) {
	        return new ProductSearchResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.mavt = source["mavt"];
	        this.tenvt = source["tenvt"];
	        this.dvt = source["dvt"];
	        this.giaban = source["giaban"];
	        this.loaithue = source["loaithue"];
	        this.vat_ra = source["vat_ra"];
	        this.tonkho = source["tonkho"];
	    }
	}
	export class SaleDTO {
	    header: BanleTam;
	    items: BanleTamCT[];
	
	    static createFrom(source: any = {}) {
	        return new SaleDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.header = this.convertValues(source["header"], BanleTam);
	        this.items = this.convertValues(source["items"], BanleTamCT);
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

