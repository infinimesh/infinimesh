import {Component} from '@angular/core';
import {NbLoginComponent} from '@nebular/auth';

@Component({
  selector: 'user-name-login',
  templateUrl: './user-name-login.component.html',
})
export class UserNameLoginComponent extends NbLoginComponent {
}
