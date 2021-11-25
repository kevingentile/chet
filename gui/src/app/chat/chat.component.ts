import { Component, OnInit } from '@angular/core';
import { ChatService } from '../chat.service';
import { webSocket, WebSocketSubject, WebSocketSubjectConfig } from 'rxjs/webSocket';
import { retry } from 'rxjs';

interface Message {
  id: string
}

@Component({
  selector: 'app-chat',
  templateUrl: './chat.component.html',
  styleUrls: ['./chat.component.scss']
})
export class ChatComponent implements OnInit {
  messages: string[]
  ws: WebSocketSubject<Message> | undefined
  constructor(private chatService: ChatService) {
    this.messages = []
   }

  ngOnInit(): void {
    // const cfg: WebSocketSubjectConfig = {
    //   url: "ws://localhost:8000/ws",
    //   binaryType: 'arraybuffer'
    // }
    this.ws = webSocket("ws://localhost:8000/ws")
    // this.ws.addEventListener('open', (event: Event) => {
    //   console.log(event)
    //   let msg = {
    //     "type": "connect",
    //     "id": "asdf",
    //   }

    //   this.ws?.send(JSON.stringify(msg))
    //   this.messages.push(`${event}`)
    // })

    // this.ws.addEventListener('message', (event: MessageEvent) => {
    //   console.log(event.data)
    // })

    this.ws.pipe(retry()).subscribe({
      next: (message: Message) => {
        console.log(message)
        
        this.messages.push(JSON.stringify(message))
      },
      error: (err) => {
        console.log(err)
      }
    })
    let msg: Message = {
        id: "asdf",
    }
    this.ws.next(msg)
  }
}
