import {Component, EventEmitter, OnInit, Output} from '@angular/core';

@Component({
  selector: 'delete-confirm',
  templateUrl: './delete-confirm.component.html',
  styleUrls: ['./delete-confirm.component.scss']
})
export class DeleteConfirmComponent implements OnInit {

  public confirm = false;
  @Output() deleteConfirmed = new EventEmitter<boolean>();

  constructor() { }

  confirmDelete() {
    this.deleteConfirmed.emit(true);
    this.confirm = false;
  }

  ngOnInit() {
  }

}
