import {EventEmitter, Injectable} from '@angular/core';
import {NbAuthService} from '@nebular/auth';
import {UrlProviderService} from './url-provider.service';
import {HttpClient} from '@angular/common/http';
import {ApiUtilService} from './api-util.service';
import {Observable} from 'rxjs';

@Injectable()
export class NamespaceService {

  selected;
  selectedChange: EventEmitter<any> = new EventEmitter();

  constructor(private http: HttpClient,
              private apiUtilService: ApiUtilService) {
    this.http = http;
    this.apiUtilService = apiUtilService;
    const selectedFromLocalStorage = localStorage.getItem('namespace');
    if (selectedFromLocalStorage) {
      this.selected = JSON.parse(selectedFromLocalStorage);
    }
  }


  getAll(): Observable<any> {
    return this.http.get(`${this.apiUtilService.getApiUrl()}/namespaces`, this.apiUtilService.getHttpOptions())
      .map((response: any) => response.namespaces);
  }

  getSelected() {
    return this.selected;
  }

  setSelected(selected) {
    this.selected = selected;
    localStorage.setItem('namespace', JSON.stringify(this.selected));
    this.selectedChange.emit(this.selected);
  }
}
