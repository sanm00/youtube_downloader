export namespace config {
	
	export class Config {
	    downloadDir: string;
	    maxConcurrent: number;
	    retryTime: number;
	
	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.downloadDir = source["downloadDir"];
	        this.maxConcurrent = source["maxConcurrent"];
	        this.retryTime = source["retryTime"];
	    }
	}

}

export namespace downloader {
	
	export class VideoInfo {
	    id: string;
	    title: string;
	    filePath: string;
	    status: string;
	    progress: number;
	    message: string;
	    // Go type: time
	    createdAt: any;
	    // Go type: time
	    completedAt?: any;
	    size: number;
	    exists: boolean;
	
	    static createFrom(source: any = {}) {
	        return new VideoInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.title = source["title"];
	        this.filePath = source["filePath"];
	        this.status = source["status"];
	        this.progress = source["progress"];
	        this.message = source["message"];
	        this.createdAt = this.convertValues(source["createdAt"], null);
	        this.completedAt = this.convertValues(source["completedAt"], null);
	        this.size = source["size"];
	        this.exists = source["exists"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
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

