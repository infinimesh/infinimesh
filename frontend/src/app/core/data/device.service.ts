import {Observable} from 'rxjs';
import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import 'rxjs/add/operator/map';
import 'rxjs/add/operator/filter';
import 'rxjs/add/operator/mergeMap';
import {ApiUtilService} from './api-util.service';
import {NamespaceService} from './namespace.service';

@Injectable()
export class DeviceService {

  constructor(private http: HttpClient,
              private apiUtilService: ApiUtilService,
              private namespaceService: NamespaceService) {

  }

  getAll(): Observable<any> {
    const namespace = this.namespaceService.getSelected();
    const url = `${this.apiUtilService.getApiUrl()}/devices?namespace=${namespace.name}`;
    return this.http.get(url, this.apiUtilService.getHttpOptions())
      .map((response: any) => response.devices);
  }

  getOne(deviceId): Observable<any> {
    const url = `${this.apiUtilService.getApiUrl()}/devices/${deviceId}`;
    return this.http.get(url, this.apiUtilService.getHttpOptions())
      .map((response: any) => response.device);
  }

  getState(deviceId): Observable<any> {
    const url = `${this.apiUtilService.getApiUrl()}/devices/${deviceId}/state`;
    return this.http.get(url, this.apiUtilService.getHttpOptions())
      .map((response: any) => response.shadow);
  }

  streamState(deviceId): Observable<any> {
    return new Observable(observer => {
      const url = `${this.apiUtilService.getApiUrl()}/devices/${deviceId}/state/stream`;
      const wsUrl = url.replace('http', 'ws')
      const ws = new WebSocket(wsUrl, ["Bearer", this.apiUtilService.getToken()]);
      ws.onmessage = function (message) {
        observer.next(JSON.parse(message.data));
      }
      return {
        unsubscribe() {
          ws.close();
        }
      };
    }).map((response: any) => response.result);
  }

  create(device): Observable<any> {
    const url = `${this.apiUtilService.getApiUrl()}/devices`;
    return this.http.post(url, {
      'device': device
    }, this.apiUtilService.getHttpOptions())
      .map((response: any) => response.device);
  }

}
