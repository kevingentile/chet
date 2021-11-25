import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http'
@Injectable({
  providedIn: 'root'
})
export class ChatService {
  constructor(private client: HttpClient) { }

  connect(): WebSocket {
    let ws = new WebSocket("localhost:8000/ws")
    return ws
  }
}
