import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { AppConfigService } from '../app-config.service';

export interface ISetPayload {
  Data: ISet[]
}

export interface ISet {
  SetId: number  // 1
  Name: string   // Genesis
  SetSeries: number
}

@Injectable({
  providedIn: 'root'
})
export class SetsDatabaseService {

  private cachedPayload: ISetPayload | undefined;

  constructor(
    private config: AppConfigService,
    private http: HttpClient) { }

  async getSetList(): Promise<ISet[]> {
    await this.ensurePayload();
    return this.cachedPayload!.Data
  }

  protected async ensurePayload() : Promise<void> {
    if (!this.cachedPayload) {
      this.cachedPayload = await this.http.get<any>(this.config.getSetsDataUrl()).toPromise();
    }
  }

  protected setPayload(cachedPayload: ISetPayload): void {
    this.cachedPayload = cachedPayload
  }

}
