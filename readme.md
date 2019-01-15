This greet service caters all the existing offerings
of the grpc which leverages the HTTP2 protocol

1.Unary
2.Server side streaming


3.Client side streaming

    usecase: 
     1. when client needs to send a lot of data(big data)
     2. When server processing is expensive and should happen
        as clients send data.
     3. When client sends data to teh server without expecting
        any response from it.   

4.Bidirectional Streaming

    usecases:
    1. When the client and server needs to send a lot of data 
       asynchronously.
    2. "Chat" protocol
    3. For a very long running message.  