package main

import (
    "fmt"
    "log"
    "net/http"
    "strings"
    "encoding/json"
    "strconv"

    "github.com/gin-gonic/gin"
    "github.com/gorilla/websocket"
)

func wsAgent(c *gin.Context) {
    _, err := GetAgentByIP(strings.Split(c.Request.RemoteAddr, ":")[0])
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    agentClients[conn] = true
    go handleAgentSocket(conn)
}

func wsWeb(c *gin.Context) {
    conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    webClients[conn] = true
    go handleWebSocket(conn)
}

func handleWebSocket(conn *websocket.Conn) {
    for {
        _, msg, err := conn.ReadMessage()
        if err != nil {
            log.Println(err)
            conn.Close()
            delete(webClients, conn)
            break
        }

        jsonData := make(map[string]interface{})
        err = json.Unmarshal(msg, &jsonData)
        if err != nil {
            log.Println(err)
            return
        }

        agentChan <- string(msg)
    }
}

func handleAgentSocket(conn *websocket.Conn) {
    for {
        _, msg, err := conn.ReadMessage()
        if err != nil {
            log.Println(err)

            deadClient := strings.Split(conn.NetConn().RemoteAddr().String(), ":")[0]
            deadAgent, err := GetAgentByIP(deadClient)
            if err != nil {
                log.Println(err)
                return
            }

            jsonData := []byte(fmt.Sprintf(`{"ID": "%s", "Status": "Dead"}`, deadAgent.ID))
            statusChan <- string(jsonData)

            conn.Close()
            delete(agentClients, conn)
            break
        }

        jsonData := make(map[string]interface{})
        err = json.Unmarshal(msg, &jsonData)
        if err != nil {
            log.Println(err)
            return
        }

        switch jsonData["OpCode"].(float64) {
        case 0:
            jsonData := []byte(fmt.Sprintf(`{"ID": "%s", "Status": "Alive"}`, strconv.FormatFloat(jsonData["ID"].(float64), 'f', -1, 64)))
            statusChan <- string(jsonData)
        case 1:
            output := jsonData["Output"].(string)
            output = strings.ReplaceAll(output, `\`, `\\`)
            jsonData := []byte(fmt.Sprintf(`{"ID": "%s", "Output": "%s"}`, strconv.FormatFloat(jsonData["ID"].(float64), 'f', -1, 64), output))
            statusChan <- string(jsonData)
        default:
            log.Println("Unknown OpCode")
        }
    }
}

func handleMsg() {
    for {
        select {
        case msg := <-statusChan:
            for client := range webClients {
                client.WriteMessage(websocket.TextMessage, []byte(msg))
            }
        case msg := <-agentChan:
            jsonData := make(map[string]interface{})
            err := json.Unmarshal([]byte(msg), &jsonData)
            if err != nil {
                log.Println(err)
                return
            }
            ip := jsonData["IP"].(string)
            for client := range agentClients {
                clientIP := strings.Split(client.NetConn().RemoteAddr().String(), ":")[0]
                if ip == clientIP {
                    client.WriteMessage(websocket.TextMessage, []byte(msg))
                }
            }
        }
    }
}
