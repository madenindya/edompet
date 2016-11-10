// // -1 -> user not exist
// func GetTotalSaldo(id string) (int64, error) {
//     var err error
//     total := int64(0)

//     if u_saldo.Id != id {
//         u_saldo, err = getUser(id)
//     }
//     if err != nil {
//         // user not exist
//     } else {
//         // add saldo if user exist
//         total = total + u_saldo.Nilai
//     }

//     // get from all cabang
//     ips := NsKelompok
//     log.Println("[CHECK] ns", ips)
//     for _, ip := range ips {
//         // url: ip/ewallet/getSaldo/user_id
//         url := fmt.Sprintf("http://%v.sisdis.ui.ac.id/ewallet/getSaldo/%v", ip, id)
//         log.Println(url)
//         resp, _ := http.Get(url)
//         defer resp.Body.Close()
//         body, _ := ioutil.ReadAll(resp.Body)
//         var sld Saldo
//         _ = json.Unmarshal(body, &sld)
//         log.Println("[CHECK]", ip, " saldo", sld.Nilai)
//         if sld.Nilai > 0 {
//             total = total + sld.Nilai
//         }
//     }
//     return total, nil
// }



// func getAllIp() ([]string, error) {
//     ips := make([]string, 0)
//     query := `
//     select ip_domisili
//     from usaldo`
//     err := db_main.Select(&ips, query)
//     return ips, err
// }


//
//
//
//
// OLD
// Handle Saldo
// resp, _ := http.Get(urlstr)
//
// TRANSFER
// res, err := http.PostForm(urlstr, url.Values{"user_id": {id}, "nilai": {nilai_str}})
//      if err != nil {
//          log.Println("[ERROR] Handler HandleTransfer PostForm", err)
//      }
