import {
  Component,
  ComponentFactory,
  ComponentFactoryResolver,
  ComponentRef,
  OnInit,
  TemplateRef,
  ViewChild,
  ViewContainerRef
} from '@angular/core';
import {DeviceService} from '../../core/data/device.service';
import {Observable} from 'rxjs';
import {NamespaceService} from '../../core/data/namespace.service';
import {NbDialogService} from "@nebular/theme";
import {SidebarComponent} from "../../theme/components/sidebar/sidebar.component";
import {DeviceEditorComponent} from '../device-editor/device-editor.component';

@Component({
  selector: 'device-registry',
  templateUrl: './device-registry.component.html',
  styleUrls: ['./device-registry.component.scss']
})
export class DeviceRegistryComponent implements OnInit {

  public devices$: Observable<any>;
  public namespaces$: Observable<any>;
  @ViewChild('sidebar') sidebar: SidebarComponent;
  @ViewChild("editorContainer", { read: ViewContainerRef }) editorContainer;
  private editorComponentRef: ComponentRef<DeviceEditorComponent>;

  constructor(private deviceService: DeviceService,
              private resolver: ComponentFactoryResolver,
              public namespaceService: NamespaceService) {
    this.namespaceService.selectedChange.subscribe(() => this.requestDevices());
  }

  ngOnInit() {
    this.requestDevices();
    this.namespaces$ = this.namespaceService.getAll();
  }

  addDevice() {
    this.editorContainer.clear();
    const factory = this.resolver.resolveComponentFactory(DeviceEditorComponent);
    this.editorComponentRef = this.editorContainer.createComponent(factory);
    this.editorComponentRef.instance.closed.subscribe(saved => {
      this.sidebar.open = false;
      this.editorComponentRef.destroy();
      if(saved) {
        this.requestDevices();
      }
    })
    this.sidebar.open = true;
  }

  delete(deviceId) {
    this.deviceService.remove(deviceId).subscribe(result => {
      this.requestDevices();
    });
  }

  saveDevice() {
    console.log("device saved");
  }

  requestDevices() {
    if (this.namespaceService.getSelected() !== undefined) {
      this.devices$ = this.deviceService.getAll();
    }
  }
}
