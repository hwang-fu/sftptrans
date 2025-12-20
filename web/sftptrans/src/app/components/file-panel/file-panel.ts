import { Component, Input, Output, EventEmitter } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FileEntry } from '../../services/api.service';

@Component({
  selector: 'app-file-panel',
  standalone: true,
  imports: [CommonModule],
  template: `
    <div class="panel">
      <div class="panel-header">
        <strong>{{ title }}</strong>
      </div>
      <div class="path-bar">
        Path: <span class="path">{{ currentPath }}</span>
      </div>
      <div class="file-list">
        <div class="file-item parent" (dblclick)="goUp()">
          [..] Parent Directory
        </div>
        <div
          *ngFor="let file of files"
          class="file-item"
          [class.selected]="selectedFile?.path === file.path"
          [class.directory]="file.isDir"
          (click)="onSelect(file)"
          (dblclick)="onOpen(file)"
        >
          <span class="icon">{{ file.isDir ? '[D]' : '[F]' }}</span>
          <span class="name">{{ file.name }}{{ file.isDir ? '/' : '' }}</span>
          <span class="size" *ngIf="!file.isDir">{{ formatSize(file.size) }}</span>
          <span class="perms">{{ file.permissions }}</span>
        </div>
      </div>
      <ng-content></ng-content>
    </div>
  `,
  styles: [`
    .panel {
      flex: 1;
      display: flex;
      flex-direction: column;
      border: 1px solid #bdc3c7;
      margin: 8px;
      background: white;
    }
    .panel-header {
      padding: 8px 12px;
      background: #3498db;
      color: white;
    }
    .path-bar {
      padding: 6px 12px;
      background: #ecf0f1;
      font-family: monospace;
      font-size: 12px;
      border-bottom: 1px solid #bdc3c7;
    }
    .path {
      color: #2c3e50;
    }
    .file-list {
      flex: 1;
      overflow-y: auto;
      font-family: 'Consolas', 'Monaco', monospace;
      font-size: 13px;
    }
    .file-item {
      display: flex;
      padding: 4px 12px;
      cursor: pointer;
      border-bottom: 1px solid #ecf0f1;
      gap: 8px;
    }
    .file-item:hover {
      background: #ebf5fb;
    }
    .file-item.selected {
      background: #3498db;
      color: white;
    }
    .file-item.parent {
      color: #7f8c8d;
      font-style: italic;
    }
    .icon {
      color: #7f8c8d;
      width: 24px;
    }
    .file-item.selected .icon {
      color: white;
    }
    .file-item.directory .name {
      color: #2980b9;
      font-weight: bold;
    }
    .file-item.selected.directory .name {
      color: white;
    }
    .name {
      flex: 1;
    }
    .size {
      color: #7f8c8d;
      width: 80px;
      text-align: right;
    }
    .file-item.selected .size {
      color: #ecf0f1;
    }
    .perms {
      color: #95a5a6;
      width: 100px;
      font-size: 11px;
    }
    .file-item.selected .perms {
      color: #ecf0f1;
    }
  `]
})
export class FilePanel {
  @Input() title = '';
  @Input() files: FileEntry[] = [];
  @Input() currentPath = '';
  @Input() selectedFile: FileEntry | null = null;

  @Output() pathChange = new EventEmitter<string>();
  @Output() fileSelect = new EventEmitter<FileEntry>();

  onSelect(file: FileEntry) {
    this.fileSelect.emit(file);
  }

  onOpen(file: FileEntry) {
    if (file.isDir) {
      this.pathChange.emit(file.path);
    }
  }

  goUp() {
    const parts = this.currentPath.split('/').filter(p => p);
    parts.pop();
    const newPath = '/' + parts.join('/');
    this.pathChange.emit(newPath || '/');
  }

  formatSize(bytes: number): string {
    if (bytes < 1024) return bytes + ' B';
    if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB';
    if (bytes < 1024 * 1024 * 1024) return (bytes / (1024 * 1024)).toFixed(1) + ' MB';
    return (bytes / (1024 * 1024 * 1024)).toFixed(1) + ' GB';
  }
}
