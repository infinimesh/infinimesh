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

      let xhr = new XMLHttpRequest();
      xhr.open(
        "GET",
        url,
        true
      );
      xhr.setRequestHeader('Authorization', 'Bearer ' + this.apiUtilService.getToken());
      xhr.onprogress = () => {
        let messages = [];

        let jsonObjects = xhr.responseText.replace(/\n$/, "").split(/\n/);
        jsonObjects.forEach(obj => {
          try {
            messages.splice(0, 0, JSON.parse(obj));
          } catch (e) {
            // ignore json parsing errors
          }
        });
        observer.next(messages[0]);
      };
      xhr.send();

      return {
        unsubscribe() {
          xhr.abort();
        }
      };
    }).map((response: any) => response.result);
  }

}
