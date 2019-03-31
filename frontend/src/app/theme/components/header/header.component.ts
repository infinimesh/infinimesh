import {Component, Input, OnInit} from '@angular/core';

import {NbMenuService, NbSidebarService} from '@nebular/theme';
import {UserService} from '../../../core/data/user.service';
import {AnalyticsService} from '../../../core/utils/analytics.service';
import {Router} from '@angular/router';
import {filter} from 'rxjs/operators';
import {NamespaceService} from '../../../core/data/namespace.service';

@Component({
  selector: 'ngx-header',
  styleUrls: ['./header.component.scss'],
  templateUrl: './header.component.html',
})
export class HeaderComponent implements OnInit {

  @Input() position = 'normal';

  user: any;
  namespaces: any;
  selectedNamespace: object;

  userMenu = [{title: 'Log out'}];

  constructor(private sidebarService: NbSidebarService,
              private menuService: NbMenuService,
              private userService: UserService,
              private router: Router,
              private namespaceService: NamespaceService,
              private analyticsService: AnalyticsService) {
  }

  ngOnInit() {
    this.menuService.onItemClick()
      .pipe(
        filter(item => item.tag === 'header-menu'),
      )
      .subscribe(() => {
        localStorage.removeItem('auth_app_token');
        this.router.navigate(['auth/login']);
      });
    this.userService.getSelf()
      .subscribe((user: any) => {
        this.user = user;
      });
    this.namespaceService.getAll().subscribe(namespaces => {
      this.namespaces = namespaces;
      const selected = this.namespaceService.getSelected();
      if (selected) {
        setTimeout(() => {
          this.selectedNamespace = this.namespaces.find(namespace => namespace.id === selected.id);
        });
      }
    });
  }

  namespaceSelectionChanged(namespace) {
    this.namespaceService.setSelected(namespace);
  }

  toggleSidebar(): boolean {
    this.sidebarService.toggle(true, 'menu-sidebar');
    return false;
  }

  goToHome() {
    this.menuService.navigateHome();
  }

}
