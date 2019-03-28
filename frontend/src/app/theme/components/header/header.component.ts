import {Component, Input, OnInit} from '@angular/core';

import {NbMenuService, NbSidebarService} from '@nebular/theme';
import {UserService} from '../../../core/data/user.service';
import {AnalyticsService} from '../../../core/utils/analytics.service';
import {Router} from "@angular/router";
import {filter} from 'rxjs/operators';

@Component({
  selector: 'ngx-header',
  styleUrls: ['./header.component.scss'],
  templateUrl: './header.component.html',
})
export class HeaderComponent implements OnInit {

  @Input() position = 'normal';

  user: any;

  userMenu = [{title: 'Log out'}];

  constructor(private sidebarService: NbSidebarService,
              private menuService: NbMenuService,
              private userService: UserService,
              private router: Router,
              private analyticsService: AnalyticsService) {
  }

  ngOnInit() {
    this.menuService.onItemClick()
      .pipe(
        filter(item => item.tag == 'header-menu'),
      )
      .subscribe(() => {
        localStorage.removeItem('auth_app_token');
        this.router.navigate(['auth/login']);
      });
    this.userService.getSelf()
      .subscribe((user: any) => {
        this.user = user;
      });
  }

  toggleSidebar(): boolean {
    this.sidebarService.toggle(true, 'menu-sidebar');
    return false;
  }

  goToHome() {
    this.menuService.navigateHome();
  }

}
