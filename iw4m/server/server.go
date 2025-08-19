package server

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/Yallamaztar/iw4m-go/iw4m"
)

type Server struct {
	iw4m *iw4m.IW4MWrapper
}

// Create a new Server wrapper
func NewServer(iw4m *iw4m.IW4MWrapper) *Server {
	return &Server{iw4m: iw4m}
}

func (s *Server) Status() ([]ServerStatus, error) {
	res, err := s.iw4m.DoRequest("/api/status")
	if err != nil {
		return nil, err
	}

	body, err := readBody(res)
	if err != nil {
		return nil, err
	}

	var status []ServerStatus
	if err := json.Unmarshal(body, &status); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return status, nil
}

func (s *Server) Info() (*ServerInfo, error) {
	res, err := s.iw4m.DoRequest("/api/info")
	if err != nil {
		return nil, err
	}

	body, err := readBody(res)
	if err != nil {
		return nil, err
	}

	var info ServerInfo
	if err := json.Unmarshal([]byte(body), &info); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return &info, nil
}

func (s *Server) MapName() (string, error) {
	doc, err := s.getDoc("/")
	if err != nil {
		return "", err
	}

	div := doc.Find("div.col-12.align-self-center.text-center.text-lg-left.col-lg-4").First()
	if div.Length() == 0 {
		return "", fmt.Errorf("map name not found")
	}

	return strings.TrimSpace(div.Find("span").First().Text()), nil
}

func (s *Server) GameMode() (string, error) {
	doc, err := s.getDoc("/")
	if err != nil {
		return "", err
	}

	div := doc.Find("div.col-12.align-self-center.text-center.text-lg-left.col-lg-4").First()
	if div.Length() == 0 {
		return "", fmt.Errorf("map name not found")
	}

	return strings.TrimSpace(div.Find("span").Eq(2).Text()), nil
}

func (s *Server) IW4MVersion() (string, error) {
	doc, err := s.getDoc("/")
	if err != nil {
		return "", err
	}

	var version string
	doc.Find("a.sidebar-link").EachWithBreak(
		func(i int, sel *goquery.Selection) bool {
			span := sel.Find("span.text-primary").First()
			if span.Length() > 0 {
				version = strings.TrimSpace(span.Text())
				return false
			}
			return true
		})

	if version == "" {
		return "", fmt.Errorf("iw4m-admin version not found")
	}

	return version, nil
}

func (s *Server) LoggedInAs() (string, error) {
	doc, err := s.getDoc("/")
	if err != nil {
		return "", err
	}

	div := doc.Find("div.sidebar-link.font-size-12.font-weight-light")
	if div.Length() == 0 {
		return "", fmt.Errorf("username not found")
	}

	return strings.TrimSpace(div.Find("colorcode").First().Text()), nil
}

func (s *Server) Rules() ([]string, error) {
	doc, err := s.getDoc("/About")
	if err != nil {
		return nil, err
	}

	var rules []string
	doc.Find("div.card.m-0.rounded").Each(
		func(i int, card *goquery.Selection) {
			h5 := card.Find("h5.text-primary.mt-0.mb-0").First()
			if h5.Length() > 0 {
				card.Find("div.rule").Each(
					func(j int, div *goquery.Selection) {
						rule := strings.TrimSpace(div.Text())
						rule = regexp.MustCompile(`\s+`).ReplaceAllString(rule, " ")
						rules = append(rules, rule)
					},
				)
			}
		})

	return rules, nil
}

func (s *Server) Reports() ([]Report, error) {
	doc, err := s.getDoc("/Action/RecentReportsForm/")
	if err != nil {
		return nil, err
	}

	var reports []Report
	doc.Find("div.rounded.bg-very-dark-dm.bg-light-ex-lm.mt-10.mb-10.p-10").Each(
		func(i int, block *goquery.Selection) {
			timestamp := strings.TrimSpace(block.Find("div.font-weight-bold").First().Text())

			block.Find("div.font-size-12").Each(
				func(j int, entry *goquery.Selection) {
					origin := strings.TrimSpace(entry.Find("a").First().Text())
					reason := strings.TrimSpace(entry.Find("span.text-white-dm.text-black-lm").First().Text())

					var target string
					targetTag := entry.Find("span.text-highlight a").First()
					if targetTag.Length() > 0 {
						target = strings.TrimSpace(targetTag.Text())
					}

					if origin != "" || reason != "" || target != "" {
						reports = append(reports, Report{
							Origin:    origin,
							Reason:    reason,
							Target:    target,
							Timestamp: timestamp,
						})
					}
				},
			)
		})

	return reports, nil
}

func (s *Server) Help() (*Help, error) {
	doc, err := s.getDoc("/Home/Help")
	if err != nil {
		return nil, err
	}

	help := &Help{
		Sections: make(map[string]HelpSection),
	}

	doc.Find("div.command-assembly-container").Each(
		func(i int, section *goquery.Selection) {
			titleTag := section.Find("h2.content-title.mb-lg-20.mt-20").First()
			if titleTag.Length() == 0 {
				return
			}

			title := strings.TrimSpace(titleTag.Text())
			if _, exists := help.Sections[title]; !exists {
				help.Sections[title] = HelpSection{
					Title:    title,
					Commands: make(map[string]Command),
				}
			}

			commands := help.Sections[title]
			section.Find("tr.d-none.d-lg-table-row.bg-dark-dm.bg-light-lm").Each(
				func(j int, cmd *goquery.Selection) {
					cells := cmd.Find("td")
					if cells.Length() < 6 {
						return
					}

					name := strings.TrimSpace(cells.Eq(0).Text())
					alias := strings.TrimSpace(cells.Eq(1).Text())
					description := strings.TrimSpace(cells.Eq(2).Text())
					requiresTarget := strings.TrimSpace(cells.Eq(3).Text())
					syntax := strings.TrimSpace(cells.Eq(4).Text())
					minLevel := strings.TrimSpace(cells.Eq(5).Text())

					commands.Commands[name] = Command{
						Alias:          alias,
						Description:    description,
						RequiresTarget: requiresTarget,
						Syntax:         syntax,
						MinLevel:       minLevel,
					}
				})
			help.Sections[title] = commands
		})

	return help, nil
}

func (s *Server) ServerIDs() ([]ServerID, error) {
	doc, err := s.getDoc("/Console")
	if err != nil {
		return nil, err
	}

	var serverIDs []ServerID
	doc.Find("select#console_server_select option").Each(
		func(i int, option *goquery.Selection) {
			name := strings.TrimSpace(option.Text())
			id, exists := option.Attr("value")
			if !exists {
				return
			}

			serverIDs = append(serverIDs, ServerID{
				Name: name,
				ID:   id,
			})
		})

	return serverIDs, nil
}

func (s *Server) ReadChat() ([]Chat, error) {
	doc, err := s.getDoc("/")
	if err != nil {
		return nil, err
	}

	var chat []Chat
	doc.Find("div.text-truncate").Each(
		func(i int, entry *goquery.Selection) {

			var sender string
			senderTag := entry.Find("span colorcode").First()
			if senderTag.Length() > 0 {
				sender = strings.TrimSpace(senderTag.Text())
			}

			var message string
			messageTags := entry.Find("span")
			if messageTags.Length() > 1 {
				messageTag := messageTags.Eq(1).Find("colorcode").First()
				if messageTag.Length() > 0 {
					message = strings.TrimSpace(messageTag.Text())
				}
			}

			if sender != "" && message != "" {
				chat = append(chat, Chat{
					Sender:  sender,
					Message: message,
				})
			}
		})

	return chat, nil
}

func (s *Server) FindPlayer(username, xuid string, count, offset, direction int) ([]FindPlayer, error) {
	if username == "" && xuid == "" {
		return nil, fmt.Errorf("username or xuid is required")
	}

	endpoint := fmt.Sprintf(
		"/api/client/find?name=%s&xuid=%s&count=%d&offset=%d&direction=%d",
		username, xuid, count, offset, direction,
	)

	res, err := s.iw4m.DoRequest(endpoint)
	if err != nil {
		return nil, err
	}

	body, err := readBody(res)
	if err != nil {
		return nil, err
	}

	var response FindPlayerResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return response.Clients, nil
}

func (s *Server) ListPlayers() ([]Players, error) {
	doc, err := s.getDoc("/")
	if err != nil {
		return nil, err
	}

	roles := map[string]string{
		"creator":       "level-color-7.no-decoration.text-truncate.ml-5.mr-5",
		"owner":         "level-color-6.no-decoration.text-truncate.ml-5.mr-5",
		"moderator":     "level-color-5.no-decoration.text-truncate.ml-5.mr-5",
		"senioradmin":   "level-color-4.no-decoration.text-truncate.ml-5.mr-5",
		"administrator": "level-color-3.no-decoration.text-truncate.ml-5.mr-5",
		"trusted":       "level-color-2.no-decoration.text-truncate.ml-5.mr-5",
		"user":          "text-light-dm.text-dark-lm.no-decoration.text-truncate.ml-5.mr-5",
		"flagged":       "level-color-1.no-decoration.text-truncate.ml-5.mr-5",
		"banned":        "level-color--1.no-decoration.text-truncate.ml-5.mr-5",
	}

	var players []Players
	for role, class := range roles {
		doc.Find("a." + class).Each(
			func(i int, sel *goquery.Selection) {
				colorcode := sel.Find("colorcode")
				if colorcode.Length() > 0 {
					players = append(players, Players{
						Role:     role,
						Name:     strings.TrimSpace(colorcode.Text()),
						ClientId: strings.TrimSpace(sel.AttrOr("href", ""))[16:],
						URL:      strings.TrimSpace(sel.AttrOr("href", "")),
					})
				}
			})
	}

	return players, nil
}

func (s *Server) StockRoles() ([]string, error) {
	doc, err := s.getDoc("/Action/editForm/?id=2&meta=")
	if err != nil {
		return nil, err
	}

	var roles []string
	selectTag := doc.Find("select[name='level']")
	if selectTag.Length() > 0 {
		selectTag.Find("option").Each(
			func(i int, sel *goquery.Selection) {
				role, exists := sel.Attr("value")
				if exists && role != "" {
					roles = append(roles, role)
				}
			})
	}
	return roles, nil
}

func (s *Server) Roles() ([]string, error) {
	doc, err := s.getDoc("/Action/editForm/?id=2&meta=")
	if err != nil {
		return nil, err
	}

	var roles []string
	doc.Find("select option").Each(func(i int, option *goquery.Selection) {
		text := strings.TrimSpace(option.Text())
		if text != "" {
			roles = append(roles, text)
		}
	})

	return roles, nil
}

func (s *Server) RecentClients(offset int) ([]RecentClient, error) {
	doc, err := s.getDoc(fmt.Sprintf("/Action/RecentClientsForm?offset=%d&count=20", offset))
	if err != nil {
		return nil, err
	}

	var clients []RecentClient
	doc.Find("div.bg-very-dark-dm.bg-light-ex-lm.p-15.rounded.mb-10").Each(
		func(i int, entry *goquery.Selection) {
			var client RecentClient

			user := entry.Find("div.d-flex.flex-row").First()
			if user.Length() > 0 {
				nameTag := user.Find("a.h4.mr-auto colorcode").First()
				if nameTag.Length() > 0 {
					client.Name = strings.TrimSpace(nameTag.Text())
				}

				linkTag := user.Find("a").First()
				if linkTag.Length() > 0 {
					client.Link = strings.TrimSpace(linkTag.AttrOr("href", ""))
				}

				tooltip := user.Find("div[data-toggle='tooltip']").First()
				if tooltip.Length() > 0 {
					client.Country = strings.TrimSpace(tooltip.AttrOr("data-title", ""))
				}
			}

			ip := entry.Find("div.align-self-center.mr-auto").First()
			if ip.Length() > 0 {
				client.IPAddress = strings.TrimSpace(ip.Text())
			}

			lastSeen := entry.Find("div.align-self-center.text-muted.font-size-12").First()
			if lastSeen.Length() > 0 {
				client.LastSeen = strings.TrimSpace(lastSeen.Text())
			}
			clients = append(clients, client)
		})

	return clients, nil
}

func (s *Server) RecentAuditLog() (*AuditLog, error) {
	doc, err := s.getDoc("/Admin/AuditLog")
	if err != nil {
		return nil, err
	}

	tbody := doc.Find("#audit_log_table_body")
	if tbody.Length() == 0 {
		return nil, nil // nothing found
	}

	tr := tbody.Find("tr.d-none.d-lg-table-row.bg-dark-dm.bg-light-lm").First()
	if tr.Length() == 0 {
		return nil, nil // no matching row
	}

	tds := tr.Find("td")
	if tds.Length() < 6 {
		return nil, fmt.Errorf("unexpected number of columns in audit log row")
	}

	originElem := tds.Eq(1).Find("a")
	targetElem := tds.Eq(2).Find("a")

	var target string
	if targetElem.Length() > 0 {
		target = strings.TrimSpace(targetElem.Text())
	} else {
		target = strings.TrimSpace(tds.Eq(2).Text())
	}

	auditLog := &AuditLog{
		Type:   strings.TrimSpace(tds.Eq(0).Text()),
		Origin: strings.TrimSpace(originElem.Text()),
		Href:   strings.TrimSpace(originElem.AttrOr("href", "")),
		Target: target,
		Data:   strings.TrimSpace(tds.Eq(4).Text()),
		Time:   strings.TrimSpace(tds.Eq(5).Text()),
	}

	return auditLog, nil
}

func (s *Server) AuditLogs(count int) ([]AuditLog, error) {
	if count <= 0 {
		count = 15
	}

	doc, err := s.getDoc("/Admin/AuditLog")
	if err != nil {
		return nil, err
	}

	tbody := doc.Find("#audit_log_table_body")
	if tbody.Length() == 0 {
		return []AuditLog{}, nil
	}

	auditLogs := []AuditLog{}
	tbody.Find("tr.d-none.d-lg-table-row.bg-dark-dm.bg-light-lm").EachWithBreak(
		func(i int, tr *goquery.Selection) bool {
			if i >= count {
				return false
			}

			tds := tr.Find("td")
			if tds.Length() < 6 {
				return true
			}

			originElem := tds.Eq(1).Find("a")
			targetElem := tds.Eq(2).Find("a")

			var target string
			if targetElem.Length() > 0 {
				target = strings.TrimSpace(targetElem.Text())
			} else {
				target = strings.TrimSpace(tds.Eq(2).Text())
			}

			auditLog := AuditLog{
				Type:   strings.TrimSpace(tds.Eq(0).Text()),
				Origin: strings.TrimSpace(originElem.Text()),
				Href:   strings.TrimSpace(originElem.AttrOr("href", "")),
				Target: target,
				Data:   strings.TrimSpace(tds.Eq(4).Text()),
				Time:   strings.TrimSpace(tds.Eq(5).Text()),
			}

			auditLogs = append(auditLogs, auditLog)
			return true
		})
	return auditLogs, nil
}

func (s *Server) Admins(role string, count int) ([]Admin, error) {
	if role == "" {
		role = "all"
	}

	doc, err := s.getDoc("/Client/Privileged")
	if err != nil {
		return nil, err
	}

	var admins []Admin
	doc.Find("table.table.mb-20").EachWithBreak(
		func(i int, table *goquery.Selection) bool {
			if count > 0 && len(admins) >= count {
				return false
			}

			headerThs := table.Find("thead tr th")
			if headerThs.Length() == 0 {
				return true
			}

			tableRole := strings.TrimSpace(headerThs.Eq(0).Text())
			if role != "all" && strings.EqualFold(strings.ToLower(tableRole), strings.ToLower(role)) {
				return true
			}

			table.Find("tbody tr").Each(
				func(j int, row *goquery.Selection) {
					if count > 0 && len(admins) >= count {
						return
					}

					name := strings.TrimSpace(row.Find("a.text-force-break").Text())
					game := "N/A"
					if badge := row.Find("div.badge"); badge.Length() > 0 {
						game = strings.TrimSpace(badge.Text())
					}
					tds := row.Find("td")
					lastConnected := strings.TrimSpace(tds.Eq(tds.Length() - 1).Text())

					admins = append(admins, Admin{
						Name:          name,
						Role:          tableRole,
						Game:          game,
						LastConnected: lastConnected,
					})
				})
			return true
		})

	return admins, nil
}

func (s *Server) TopPlayers(count int) ([]TopPlayer, error) {
	doc, err := s.getDoc(fmt.Sprintf("/Stats/GetTopPlayersAsync?offset=0&count=%d&serverId=0", count))
	if err != nil {
		return nil, err
	}

	var players []TopPlayer
	doc.Find("div.card.m-0.mt-15.p-20.d-flex.flex-column.flex-md-row.justify-content-between").Each(
		func(i int, entry *goquery.Selection) {
			rankDiv := entry.Find("div.d-flex.flex-column.w-full.w-md-quarter")
			if rankDiv.Length() == 0 {
				return
			}

			rank := strings.TrimSpace(rankDiv.Find("div.d-flex.text-muted div").First().Text())
			player := TopPlayer{
				Rank:  "#" + rank,
				Stats: map[string]string{},
			}

			nameTag := rankDiv.Find("div.d-flex.flex-row")
			if nameTag.Length() > 0 {
				colorcode := nameTag.Find("colorcode").First()
				if colorcode.Length() > 0 {
					player.Name = strings.TrimSpace(colorcode.Text())
				}
				if link, exists := nameTag.Find("a").Attr("href"); exists {
					player.Link = link
				}
			}

			rating := strings.TrimSpace(rankDiv.Find("div.font-size-14 span").First().Text())
			player.Rating = rating

			statsTag := rankDiv.Find("div.d-flex.flex-column.font-size-12.text-right.text-md-left")
			statsTag.Find("div").Each(func(j int, div *goquery.Selection) {
				primary := strings.TrimSpace(div.Find("span.text-primary").Text())
				secondary := strings.TrimSpace(div.Find("span.text-muted").Text())
				if primary != "" && secondary != "" {
					player.Stats[secondary] = primary
				}
			})

			players = append(players, player)
		})

	return players, nil
}

func (s *Server) PlayerCount() int {
	players, err := s.ListPlayers()
	if err != nil {
		return 0
	}
	return len(players)
}

func (s *Server) IsServerFull() bool {
	info, _ := s.Info()
	return info.TotalConnectedClients >= info.TotalClientSlots
}

func (s *Server) FindAdmin(username string) Admin {
	admins, _ := s.Admins("all", 1000)
	for _, admin := range admins {
		if strings.EqualFold(strings.ToLower(admin.Name), strings.ToLower(username)) {
			return admin
		}
	}
	return Admin{}
}
