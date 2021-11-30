import { Component, OnInit } from '@angular/core';
import { ChatService } from '../chat.service';
import { webSocket, WebSocketSubject, WebSocketSubjectConfig } from 'rxjs/webSocket';
import { catchError, delay, Observable, retry, retryWhen, take } from 'rxjs';

interface Message {
  id: string
  timestamp: number
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
    this.ws = webSocket("ws://localhost:8000/ws?roomId=aef24add-ea80-437b-a030-76d6e993ae52")

    this.ws.pipe(
      retryWhen(
        (errors: Observable<any>) => errors.pipe(
          delay(5000),
          take(10)
        )
      ),
    ).subscribe({
      next: (message: Message) => {
        console.log(message)

        this.messages.push(JSON.stringify(message))
        if (this.messages.length > 15) {
          this.messages = this.messages.slice(1)
        }
      },
      error: (err) => {
        // console.log(err)
      }
    })
  }
}
