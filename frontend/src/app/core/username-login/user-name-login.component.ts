import {Component, OnInit} from '@angular/core';
import {NbLoginComponent} from '@nebular/auth';
import {environment} from '../../../environments/environment';

const apiServerUrl = environment.apiServerUrl;

@Component({
  selector: 'user-name-login',
  templateUrl: './user-name-login.component.html',
})
export class UserNameLoginComponent extends NbLoginComponent implements OnInit {

  public showRegisterLink = false;

  ngOnInit(): void {
    if (apiServerUrl.includes('api.infinimesh.io')) {
      this.showRegisterLink = true;
    }
  }
}
