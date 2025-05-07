# TechFinder
TechFinder, verilen bir domain listesi üzerinden web sitelerinin kullandığı teknolojileri tespit etmek amacıyla geliştirilmiş bir analiz aracıdır.

### Nasıl Çalışır?
TechFinder, hedef web sitelerine yaptığı analizlerde aşağıdaki yöntemleri kullanarak kullanılan teknolojileri belirlemeye çalışır:

- HTTP Response Header bilgileri üzerinden sunucu, framework veya diğer teknolojiler belirlenmeye çalışılır.

- HTML Body içeriğinden:

    - Yorum satırları taranarak kullanılan teknolojilerin adları ve sürüm numaraları bulunmaya çalışılır.

    - CSS ve JavaScript dosyalarına ait referanslar analiz edilerek regex ile teknoloji tespiti yapılır.

- Cookie değerleri üzerinden teknolojilere dair ipuçları toplanır. 
    - Web uygulamasında kullanılan altyapı hakkında ipuçları elde edilir.

- Elde ettiği bileşen bilgilerini <b>[Snyk](https://snyk.io/)</b> ve <b>[CVEDetails](https://www.cvedetails.com/)</b> adreslerinden kontrol eder.

- Ayrıca, hedefe doğrudan analiz yapılmadan önce:
    -   curl komutu ile hedef sitenin canlı olup olmadığı, HTTP yönlendirmesi (3xx) yapıp yapmadığı kontrol edilir ve yönlendirme takip edilir.

----

<p style="font-size: smaller; color: gray;"><em>Araç yanlışlık payı içerebilir manuel teyit edilmelidir.</em></p>
