-- create news

CREATE TABLE IF NOT EXISTS news_types (
    type_id SERIAL PRIMARY KEY,
    type_name VARCHAR(50) NOT NULL
);

CREATE TABLE IF NOT EXISTS news (
    news_id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    type_id INT NOT NULL,
    detail_url TEXT NULL,
    cover_image TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (type_id) REFERENCES news_types(type_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS news_images (
    image_id SERIAL PRIMARY KEY,
    news_id INT NOT NULL,
    file_image TEXT NOT NULL,
    FOREIGN KEY (news_id) REFERENCES news(news_id) ON DELETE CASCADE
);

-- insert news

INSERT INTO news_types(type_name) VALUES
('ข่าวประชาสัมพันธ์'),
('ทุนการศึกษา'),
('รางวัลที่ได้รับ'),
('กิจกรรมของภาควิชา');

INSERT INTO news(title,content,type_id,detail_url,cover_image) VALUES
('คู่มือแนะนำนักศึกษาใหม่ ปีการศึกษา 2568',
    'คู่มือแนะนำ การลงทะเบียนรายวิชาเรียน การเพิ่ม - ถอน - การเปลี่ยนกลุ่มเรียน ในรายวิชาเดิม การถอนรายวิชาโดยติดสัญลักษณ์ W ช่องทางการชำระค่าธรรมเนียมการศึกษา เอกสารแนบเบิกค่าเล่าเรียน การลาพักการศึกษา นักศึกษาที่คาดว่าจะสำเร็จการศึกษา การขอหนังสือสำคัญทางการศึกษา การติดต่อเจ้าหน้าที่งานทะเบียนและประมวลผล ช่องทางติดต่อคณะวิชา กองกิจการนักศึกษา การให้บริการสำหรับนักศึกษา วิทยาเขตสารสนเทศเพชรบุรี สำนักดิจิทัลเทคโนโลยี หอสมุด มหาวิทยาลัยศิลปากร ประกันอุบัติเหตุส่วนบุคคล ศูนย์บริการสุขภาพ คณะเภสัชศาสตร์ บริการด้านสุขภาพร่างกาย เครื่องแบบนักศึกษา / บัตรนักศึกษาอิเล็กทรอนิกส์',
    (SELECT type_id FROM news_types WHERE type_name = 'ข่าวประชาสัมพันธ์' LIMIT 1),
    'https://drive.google.com/file/d/1FnJRketlluku27HQqs-mVPnxMJ5wFxqZ/view?usp=sharing&fbclid=IwZXh0bgNhZW0CMTAAAR3M6TbS4tN1DUSa-z2NaM1ekAOELZp7TnsIhJWC6g_dvfz-sD_b0La0S7U_aem_K5_Zhi2eLKhIpJxaeszOlQ',
    'http://localhost:9000/images/news/manual.jpg'),
('เปิดรับสมัครปริญญาโทและปริญญาเอกสาขาเทคโนโลยีสารสนเทศและนวัตกรรมดิจิทัล',
    'เปิดรับสมัครแล้ว! ภาควิชาคอมพิวเตอร์ คณะวิทยาศาสตร์ มหาวิทยาลัยศิลปากร เปิดรับสมัครนักศึกษาระดับบัณฑิตศึกษา ปริญญาโท - ปริญญาเอก หลักสูตร IT สำหรับคนทำงาน สู่ผู้เชี่ยวชาญด้านวิจัยและนวัตกรรม เน้นเรียนรู้กระบวนทำวิจัยเพื่อแก้ปัญหาจริง มีทุนผู้ช่วยวิจัยและทุนนำเสนอผลงานในงานประชุมวิชาการ มีทั้งภาคปกติ (เรียนวันธรรมดา) และโครงการพิเศษ (เรียน ส อา) มี 3 แผนการเรียนให้เลือก Data Science, Project Management, DevOps เปิดรับสมัครรอบ 1 วันที่ 16 ธ.ค. 67 - 21 ก.พ. 68 เปิดรับสมัครรอบ 2 วันี่ต่ 3 มี.ค. 68 - 7 พ.ค. 68 สมัครง่าย ๆ ผ่านช่องทางออนไลน์ได้ที่ https://graduate.su.ac.th มีข้อสงสัยสอบถามเพิ่มเติมได้ที่ เพจ Facebook : https://www.facebook.com/computingsu/ หรือ โทร 034-272-923',
    (SELECT type_id FROM news_types WHERE type_name = 'ข่าวประชาสัมพันธ์' LIMIT 1),
    'https://graduate.su.ac.th',
    'http://localhost:9000/images/news/degree1.jpg'),
('เปิดรับสมัครทุนการศึกษา คณะวิทยาศาสตร์ มหาวิทยาลัยศิลปากร สำหรับนักศึกษาระดับปริญญาตรี ภาคการศึกษาต้น ปีการศึกษา 2568',
    'เปิดรับสมัครทุนการศึกษา! คณะวิทยาศาสตร์ มหาวิทยาลัยศิลปากร สำหรับนักศึกษาระดับปริญญาตรี ภาคการศึกษาต้น ปีการศึกษา 2568 เปิดรับสมัคร: ตั้งแต่วันที่ 1 กรกฎาคม - 29 สิงหาคม 2568 คุณสมบัติผู้สมัคร: 1.เป็นนักศึกษาระดับปริญญาตรี ของคณะวิทยาศาสตร์ 2.ขาดแคลนทุนทรัพย์ และมีเกรดเฉลี่ยสะสมไม่ต่ำกว่า 2.00 3.นักศึกษาที่กู้ยืม กยศ./กรอ. สามารถสมัครได้',
    (SELECT type_id FROM news_types WHERE type_name = 'ทุนการศึกษา' LIMIT 1),
    'https://drive.su.ac.th/s/tHXDdo3mpJ9L45d',
    'http://localhost:9000/images/news/scholarship.jpg'),
('ขอแสดงความยินดีกับ ผู้ช่วยศาสตราจารย์ ดร.อรวรรณ เชาวลิต อาจารย์ประจำภาควิชาคอมพิวเตอร์',
    'คณะวิทยาศาสตร์ มหาวิทยาลัยศิลปากร ขอแสดงความยินดีกับ ผู้ช่วยศาสตราจารย์ ดร.อรวรรณ เชาวลิต ภาควิชาคอมพิวเตอร์ ในโอกาสที่ตีพิมพ์ผลงานวิจัยในวารสาร ICIC Express Letters, Part B: Applications ในฐานข้อมูล Scopus (ScimagoJR, Quartile 4) เรื่อง Multivariate time series forecasting Thailand soybean meal price with deep learning models ICIC Express Letters, Part B: Applications 2025, 16(3), 317-323. http://www.icicelb.org/ellb/contents/2025/3/elb-16-03-10.pdf SDG 8 DECENT WORK AND ECONOMIC GROWTH',
    (SELECT type_id FROM news_types WHERE type_name = 'รางวัลที่ได้รับ' LIMIT 1),
    'http://www.icicelb.org/ellb/contents/2025/3/elb-16-03-10.pdf',
    'http://localhost:9000/images/news/reward-orawan.jpg'),
('ขอแสดงความยินดีกับผู้ช่วยศาสตราจารย์ ดร.ปัญญนัท อ้นพงษ์ เนื่องในโอกาสได้รับการแต่งตั้งให้ดำรงตำแหน่งทางวิชาการ',
    'ภาควิชาคอมพิวเตอร์ คณะวิทยาศาสตร์ มหาวิทยาลัยศิลปากร ขอแสดงความยินดีกับผู้ช่วยศาสตราจารย์ ดร.ปัญญนัท อ้นพงษ์ เนื่องในโอกาสได้รับการแต่งตั้งให้ดำรงตำแหน่งทางวิชาการ ตำแหน่ง “ผู้ช่วยศาสตราจารย์” ในสาขาวิชาเทคโนโลยีสารสนเทศ',
    (SELECT type_id FROM news_types WHERE type_name = 'รางวัลที่ได้รับ' LIMIT 1),
    'https://www.facebook.com/ScienceSilpakornUniversity/posts/คณะวิทยาศาสตร์-มหาวิทยาลัยศิลปากร-ขอแสดงความยินดีกับผู้ช่วยศาสตราจารย์-ดรปัญญนัท/982826210555894/',
    'http://localhost:9000/images/news/reward-panyanut.jpg'),
('ขอแสดงความยินดีกับ นายภากร กัทชลี นักศึกษา ปริญญาเอก หลักสูตรเทคโนโลยีสารสนเทศและนวัตกรรมดิจิทัล',
    'ภาควิชาคอมพิวเตอร์ ขอแสดงความยินดีกับ นายภากร กัทชลี นักศึกษา ปริญญาเอก หลักสูตรเทคโนโลยีสารสนเทศและนวัตกรรมดิจิทัล ที่ได้รับรางวัล Best Presentation ในงาน The 6th Asia Joint Conference on Computing (AJCC 2025) ณ เมือง Osaka ประเทศญี่ปุ่นจาก Paper หัวข้อ "Intention Classification of Chinese Topics on Thai Facebook Pages Using Transformer Models with Emotional Features" อาจารย์ที่ปรึกษา อ.ดร.สัจจาภรณ์ ไวจรรยา และ ผศ.ดร.ณัฐโชติ พรหมฤทธิ์',
    (SELECT type_id FROM news_types WHERE type_name = 'รางวัลที่ได้รับ' LIMIT 1),
    'https://www.facebook.com/computingsu/posts/-ภาควิชาคอมพิวเตอร์-ขอแสดงความยินดีกับ-นายภากร-กัทชลี-นักศึกษา-ปริญญาเอก-หลักสูต/122207080106177331/',
    'http://localhost:9000/images/news/reward-phakon1.jpg'),
('ขอแสดงความยินดีกับ  นายณัฐพิสิษ กังวานธรรมกุล,นายนัทธพงศ์ เป็กทองและนาย พิรพัฒน์ ยิ่งแก้ว นักศึกษาสาขาวิชาเทคโนโลยีสารสนเทศ กับผลงาน "การพัฒนาระบบอํานวยการสอบ คณะวิทยาศาสตร์ มหาวิทยาลัยศิลปากร',
    'ภาควิชาคอมพิวเตอร์ ขอแสดงความยินดีกับนักศึกษาและอาจารย์ที่ปรึกษา ที่เข้าร่วมนำเสนอผลงานและได้รับรางวัลในงาน การประชุมวิชาการระดับปริญญาตรีด้านคอมพิวเตอร์ภูมิภาคเอเชีย ครั้งที่ 13 (The 13th Asia Undergraduate Conference on Computing: AUCC 2025) ในระหว่างวันที่ 19 - 21 กุมภาพันธ์ พ.ศ. 2568 ณ มหาวิทยาลัยราชภัฏอุตรดิตถ์ จังหวัดอุตรดิตถ์ รางวัลรองชนะเลิศอันดับ 1 Track นวัตกรรม กับผลงาน "การพัฒนาระบบอํานวยการสอบ คณะวิทยาศาสตร์ มหาวิทยาลัยศิลปากร (Development of examination Management system for the Faculty of Science)" ซึ่งผลงานนี้ได้มีการนำไปใช้จริงทดแทนระบบเดิม ในการจัดสอบของคณะฯ เมื่อช่วงสอบกลางภาคการศึกษาที่ผ่านมา นักศึกษาที่เข้านำเสนอ - นาย ณัฐพิสิษ กังวานธรรมกุล - นาย นัทธพงศ์ เป็กทอง - นาย พิรพัฒน์ ยิ่งแก้ว นักศึกษาสาขาวิชาเทคโนโลยีสารสนเทศ อาจารย์ที่ปรึกษา ผศ.ดร.ณัฐโชติ พรหมฤทธิ์ และ อ.ดร.สัจจาภรณ์ ไวจรรยา',
    (SELECT type_id FROM news_types WHERE type_name = 'รางวัลที่ได้รับ' LIMIT 1),
    'https://www.facebook.com/computingsu/posts/-ภาควิชาคอมพิวเตอร์-ขอแสดงความยินดีกับนักศึกษาและอาจารย์ที่ปรึกษา-ที่เข้าร่วมนำเ/122194936622177331/',
    'http://localhost:9000/images/news/reward2nd1.jpg'),
('งานสานสัมพันธ์ภาคคอมพิวเตอร์',
    'เตรียมตัวให้พร้อม! งานสานสัมพันธ์ภาคคอมพิวเตอร์กำลังจะเริ่มขึ้น! เจอกันวันที่ 28 กุมภาพันธ์ 2568 ณ ลานจอดรถตึก 4 เวลา 17:30 - 21:30 น. ธีม: ย้อนยุค ไม่จำกัดช่วงเวลา! จะเป็นยุคหิน มนุษย์ถ้ำ อารยธรรมโบราณ หรือยุคใดก็ได้ จัดเต็มมาให้สุด! มีรางวัลแต่งกายยอดเยี่ยม 3 รางวัล!ใครแต่งตัวเข้าธีมและโดดเด่นที่สุด มีสิทธิ์คว้ารางวัลไปเลย! สนุกไปกับการแสดงสุดพิเศษ, เกมสุดมันส์, ลุ้นรางวัล และดนตรีสดปิดท้าย! ใครลงชื่อไว้แล้ว เจอกันแน่นอน! เตรียมชุดให้พร้อม แล้วมาสนุกไปด้วยกัน!',
    (SELECT type_id FROM news_types WHERE type_name = 'กิจกรรมของภาควิชา' LIMIT 1),
    '',
    'http://localhost:9000/images/news/cpsu_event.jpg');

INSERT INTO news_images(news_id,file_image) VALUES
((SELECT news_id FROM news WHERE title = 'คู่มือแนะนำนักศึกษาใหม่ ปีการศึกษา 2568' LIMIT 1),
    'http://localhost:9000/images/news/manual.jpg'),
((SELECT news_id FROM news WHERE title = 'เปิดรับสมัครปริญญาโทและปริญญาเอกสาขาเทคโนโลยีสารสนเทศและนวัตกรรมดิจิทัล' LIMIT 1),
    'http://localhost:9000/images/news/degree1.jpg'),
((SELECT news_id FROM news WHERE title = 'เปิดรับสมัครปริญญาโทและปริญญาเอกสาขาเทคโนโลยีสารสนเทศและนวัตกรรมดิจิทัล' LIMIT 1),
    'http://localhost:9000/images/news/degree2.jpg'),
((SELECT news_id FROM news WHERE title = 'เปิดรับสมัครปริญญาโทและปริญญาเอกสาขาเทคโนโลยีสารสนเทศและนวัตกรรมดิจิทัล' LIMIT 1),
    'http://localhost:9000/images/news/degree3.jpg'),
((SELECT news_id FROM news WHERE title = 'เปิดรับสมัครปริญญาโทและปริญญาเอกสาขาเทคโนโลยีสารสนเทศและนวัตกรรมดิจิทัล' LIMIT 1),
    'http://localhost:9000/images/news/degree4.jpg'),
((SELECT news_id FROM news WHERE title = 'เปิดรับสมัครทุนการศึกษา คณะวิทยาศาสตร์ มหาวิทยาลัยศิลปากร สำหรับนักศึกษาระดับปริญญาตรี ภาคการศึกษาต้น ปีการศึกษา 2568' LIMIT 1),
    'http://localhost:9000/images/news/scholarship.jpg'),
((SELECT news_id FROM news WHERE title = 'ขอแสดงความยินดีกับ ผู้ช่วยศาสตราจารย์ ดร.อรวรรณ เชาวลิต อาจารย์ประจำภาควิชาคอมพิวเตอร์' LIMIT 1),
    'http://localhost:9000/images/news/reward-orawan.jpg'),
((SELECT news_id FROM news WHERE title = 'ขอแสดงความยินดีกับผู้ช่วยศาสตราจารย์ ดร.ปัญญนัท อ้นพงษ์ เนื่องในโอกาสได้รับการแต่งตั้งให้ดำรงตำแหน่งทางวิชาการ' LIMIT 1),
    'http://localhost:9000/images/news/reward-panyanut.jpg'),
((SELECT news_id FROM news WHERE title = 'ขอแสดงความยินดีกับ นายภากร กัทชลี นักศึกษา ปริญญาเอก หลักสูตรเทคโนโลยีสารสนเทศและนวัตกรรมดิจิทัล' LIMIT 1),
    'http://localhost:9000/images/news/reward-phakon1.jpg'),
((SELECT news_id FROM news WHERE title = 'ขอแสดงความยินดีกับ นายภากร กัทชลี นักศึกษา ปริญญาเอก หลักสูตรเทคโนโลยีสารสนเทศและนวัตกรรมดิจิทัล' LIMIT 1),
    'http://localhost:9000/images/news/reward-phakon2.jpg'),
((SELECT news_id FROM news WHERE title = 'ขอแสดงความยินดีกับ  นายณัฐพิสิษ กังวานธรรมกุล,นายนัทธพงศ์ เป็กทองและนาย พิรพัฒน์ ยิ่งแก้ว นักศึกษาสาขาวิชาเทคโนโลยีสารสนเทศ กับผลงาน "การพัฒนาระบบอํานวยการสอบ คณะวิทยาศาสตร์ มหาวิทยาลัยศิลปากร' LIMIT 1),
    'http://localhost:9000/images/news/reward2nd1.jpg'),
((SELECT news_id FROM news WHERE title = 'ขอแสดงความยินดีกับ  นายณัฐพิสิษ กังวานธรรมกุล,นายนัทธพงศ์ เป็กทองและนาย พิรพัฒน์ ยิ่งแก้ว นักศึกษาสาขาวิชาเทคโนโลยีสารสนเทศ กับผลงาน "การพัฒนาระบบอํานวยการสอบ คณะวิทยาศาสตร์ มหาวิทยาลัยศิลปากร' LIMIT 1),
    'http://localhost:9000/images/news/reward2nd2.jpg'),
((SELECT news_id FROM news WHERE title = 'ขอแสดงความยินดีกับ  นายณัฐพิสิษ กังวานธรรมกุล,นายนัทธพงศ์ เป็กทองและนาย พิรพัฒน์ ยิ่งแก้ว นักศึกษาสาขาวิชาเทคโนโลยีสารสนเทศ กับผลงาน "การพัฒนาระบบอํานวยการสอบ คณะวิทยาศาสตร์ มหาวิทยาลัยศิลปากร' LIMIT 1),
    'http://localhost:9000/images/news/reward2nd3.jpg'),
((SELECT news_id FROM news WHERE title = 'ขอแสดงความยินดีกับ  นายณัฐพิสิษ กังวานธรรมกุล,นายนัทธพงศ์ เป็กทองและนาย พิรพัฒน์ ยิ่งแก้ว นักศึกษาสาขาวิชาเทคโนโลยีสารสนเทศ กับผลงาน "การพัฒนาระบบอํานวยการสอบ คณะวิทยาศาสตร์ มหาวิทยาลัยศิลปากร' LIMIT 1),
    'http://localhost:9000/images/news/reward2nd4.jpg'),
((SELECT news_id FROM news WHERE title = 'งานสานสัมพันธ์ภาคคอมพิวเตอร์' LIMIT 1),
    'http://localhost:9000/images/news/cpsu_event.jpg');

-- create courses

CREATE TABLE IF NOT EXISTS career_paths (
    career_paths_id SERIAL PRIMARY KEY,
    career_paths TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS plo (
    plo_id SERIAL PRIMARY KEY,
    plo TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS courses (
    course_id VARCHAR(10) PRIMARY KEY,
    degree VARCHAR(100) NOT NULL,
    major VARCHAR(100) NOT NULL,
    year INT NOT NULL,
    thai_course VARCHAR(255) NOT NULL,
    eng_course VARCHAR(255) NOT NULL,
    thai_degree VARCHAR(255) NOT NULL,
    eng_degree VARCHAR(255) NOT NULL,
    admission_req TEXT NOT NULL,
    graduation_req TEXT NOT NULL,
    philosophy TEXT NOT NULL,
    objective TEXT NOT NULL,
    tuition TEXT NOT NULL,
    credits VARCHAR(50) NOT NULL,
    career_paths_id INT NOT NULL,
    plo_id INT NOT NULL,
    detail_url TEXT NOT NULL,
    status VARCHAR(25) NOT NULL,
    FOREIGN KEY (career_paths_id) REFERENCES career_paths(career_paths_id) ON DELETE CASCADE,
    FOREIGN KEY (plo_id) REFERENCES plo(plo_id) ON DELETE CASCADE
);

-- insert courses

COPY career_paths(career_paths)
FROM '/docker-entrypoint-initdb.d/csv/course/career_paths.csv'
WITH (FORMAT csv, HEADER true, ENCODING 'UTF8');

COPY plo(plo)
FROM '/docker-entrypoint-initdb.d/csv/course/plo.csv'
WITH (FORMAT csv, HEADER true, ENCODING 'UTF8');

COPY courses(course_id,degree,major,year,thai_course,eng_course,thai_degree,eng_degree,admission_req,graduation_req,philosophy,objective,tuition,credits,career_paths_id,plo_id,detail_url,status)
FROM '/docker-entrypoint-initdb.d/csv/course/courses.csv'
WITH (FORMAT csv, HEADER true, ENCODING 'UTF8');

SELECT setval('career_paths_career_paths_id_seq', (SELECT MAX(career_paths_id) FROM career_paths));
SELECT setval('plo_plo_id_seq', (SELECT MAX(plo_id) FROM plo));

-- create course structure

CREATE TABLE IF NOT EXISTS course_structure (
    course_structure_id SERIAL PRIMARY KEY,
    course_id VARCHAR(10) NOT NULL,
    detail TEXT NOT NULL,
    FOREIGN KEY (course_id) REFERENCES courses(course_id) ON DELETE CASCADE
);

-- create roadmap

CREATE TABLE IF NOT EXISTS roadmap (
    roadmap_id SERIAL PRIMARY KEY,
    course_id VARCHAR(10) NOT NULL,
    roadmap_url TEXT NOT NULL,
    FOREIGN KEY (course_id) REFERENCES courses(course_id) ON DELETE CASCADE
);

-- insert roadmap

INSERT INTO roadmap(course_id,roadmap_url) VALUES
((SELECT course_id FROM courses WHERE thai_course = '(วท.บ) หลักสูตรวิทยาศาสตรบัณฑิต สาขาวิชาวิทยาการคอมพิวเตอร์ 2565' LIMIT 1),
    'http://localhost:9000/images/roadmap/roadmap_CS_65.jpg'),
((SELECT course_id FROM courses WHERE thai_course = '(วท.บ) หลักสูตรวิทยาศาสตรบัณฑิต สาขาวิชาเทคโนโลยีสารสนเทศ 2565' LIMIT 1),
    'http://localhost:9000/images/roadmap/roadmap_IT_65.jpg');

-- create subject

CREATE TABLE IF NOT EXISTS description (
    description_id VARCHAR(6) PRIMARY KEY,
    description_thai TEXT NULL,
    description_eng TEXT NULL
);

CREATE TABLE IF NOT EXISTS clo (
    clo_id VARCHAR(6) PRIMARY KEY,
    clo TEXT NULL
);

CREATE TABLE IF NOT EXISTS subjects (
    id SERIAL PRIMARY KEY,
    subject_id VARCHAR(10) NOT NULL,
    course_id VARCHAR(10) NOT NULL,
    plan_type VARCHAR(50) NOT NULL,
    semester VARCHAR(50) NOT NULL,
    thai_subject VARCHAR(100) NOT NULL,
    eng_subject VARCHAR(100) NULL,
    credits VARCHAR(50) NOT NULL,
    compulsory_subject VARCHAR(255) NULL,
    condition VARCHAR(255) NULL,
    description_id VARCHAR(6) NULL,
    clo_id VARCHAR(6) NULL,
    FOREIGN KEY (course_id) REFERENCES courses(course_id) ON DELETE CASCADE,
    FOREIGN KEY (description_id) REFERENCES description(description_id) ON DELETE CASCADE,
    FOREIGN KEY (clo_id) REFERENCES clo(clo_id) ON DELETE CASCADE
);

-- insert subject

COPY description(description_id,description_thai,description_eng)
FROM '/docker-entrypoint-initdb.d/csv/subject/BS_description.csv'
WITH (FORMAT csv, HEADER true, ENCODING 'UTF8');

COPY description(description_id,description_thai,description_eng)
FROM '/docker-entrypoint-initdb.d/csv/subject/MSIT66_description.csv'
WITH (FORMAT csv, HEADER true, ENCODING 'UTF8');

COPY description(description_id,description_thai,description_eng)
FROM '/docker-entrypoint-initdb.d/csv/subject/DSIT66_description.csv'
WITH (FORMAT csv, HEADER true, ENCODING 'UTF8');

COPY clo(clo_id,clo)
FROM '/docker-entrypoint-initdb.d/csv/subject/BS_clo.csv'
WITH (FORMAT csv, HEADER true, ENCODING 'UTF8');

COPY clo(clo_id,clo)
FROM '/docker-entrypoint-initdb.d/csv/subject/MSIT66_clo.csv'
WITH (FORMAT csv, HEADER true, ENCODING 'UTF8');

COPY clo(clo_id,clo)
FROM '/docker-entrypoint-initdb.d/csv/subject/DSIT66_clo.csv'
WITH (FORMAT csv, HEADER true, ENCODING 'UTF8');

COPY subjects(subject_id,course_id,plan_type,semester,thai_subject,eng_subject,credits,compulsory_subject,condition,description_id,clo_id)
FROM '/docker-entrypoint-initdb.d/csv/subject/BS_subjects.csv'
WITH (FORMAT csv, HEADER true, ENCODING 'UTF8');

COPY subjects(subject_id,course_id,plan_type,semester,thai_subject,eng_subject,credits,compulsory_subject,condition,description_id,clo_id)
FROM '/docker-entrypoint-initdb.d/csv/subject/MSIT66_subjects.csv'
WITH (FORMAT csv, HEADER true, ENCODING 'UTF8');

COPY subjects(subject_id,course_id,plan_type,semester,thai_subject,eng_subject,credits,compulsory_subject,condition,description_id,clo_id)
FROM '/docker-entrypoint-initdb.d/csv/subject/DSIT66_subjects.csv'
WITH (FORMAT csv, HEADER true, ENCODING 'UTF8');

SELECT setval('subjects_id_seq', (SELECT MAX(id) FROM subjects));

-- create personnel

CREATE TABLE IF NOT EXISTS department_position (
    department_position_id SERIAL PRIMARY KEY,
    department_position_name VARCHAR(50) NOT NULL
);

CREATE TABLE IF NOT EXISTS academic_position (
    academic_position_id SERIAL PRIMARY KEY,
    thai_academic_position VARCHAR(50) NOT NULL,
    eng_academic_position VARCHAR(50) NULL
);

CREATE TABLE IF NOT EXISTS personnels (
    personnel_id SERIAL PRIMARY KEY,
    type_personnel VARCHAR(50) NOT NULL,
    department_position_id INT NOT NULL,
    academic_position_id INT NULL,
    thai_name VARCHAR(50) NOT NULL,
    eng_name VARCHAR(50) NOT NULL,
    education TEXT NULL,
    related_fields TEXT NULL,
    email VARCHAR(100) NULL,
    website TEXT NULL,
    file_image TEXT NOT NULL,
    scopus_id VARCHAR(50) NULL,
    FOREIGN KEY (department_position_id) REFERENCES department_position(department_position_id) ON DELETE CASCADE,
    FOREIGN KEY (academic_position_id) REFERENCES academic_position(academic_position_id) ON DELETE CASCADE
);

-- insert personnel

INSERT INTO department_position(department_position_name) VALUES
('หัวหน้าภาควิชา'),('รองหัวหน้าภาควิชาฯ ฝ่ายบริหาร'),('รองหัวหน้าภาควิชาฯ'),('อาจารย์ประจำภาควิชา'),
('นักวิชาการอุดมศึกษาชำนาญการ'),('นักวิชาการอุดมศึกษาปฏิบัติการ'),('นักวิชาการอุดมศึกษา (ประจำหลักสูตรวิทยาการข้อมูล)'),('นักวิชาการอุดมศึกษา'),('นักเทคโนโลยีสารสนเทศ'),('นักคอมพิวเตอร์'),('พนักงานทั่วไป');

INSERT INTO academic_position(thai_academic_position,eng_academic_position) VALUES
('ศ.ดร.','Prof.Dr.'),
('ศ.','Prof.'),
('รศ.ดร.','Assoc.Prof.Dr.'),
('รศ.','Assoc.Prof.'),
('ผศ.ดร.','Asst.Prof.Dr.'), 
('ผศ.','Asst.Prof.'), 
('อ.ดร.','Dr.'), 
('ดร.','Dr.'), 
('อ.',''); 

COPY personnels(type_personnel,department_position_id,academic_position_id,thai_name,eng_name,education,related_fields,email,website,file_image,scopus_id)
FROM '/docker-entrypoint-initdb.d/csv/personnel/personnels.csv'
WITH (FORMAT csv, HEADER true, ENCODING 'UTF8');

SELECT setval('personnels_personnel_id_seq', (SELECT MAX(personnel_id) FROM personnels));

-- create research

CREATE TABLE IF NOT EXISTS research (
    research_id SERIAL PRIMARY KEY,
    personnel_id INT NOT NULL,
    title TEXT NOT NULL,
    journal VARCHAR(255) NOT NULL,
    year INT NOT NULL,
    volume VARCHAR(50),
    issue VARCHAR(50),
    pages VARCHAR(50),
    doi TEXT,
    cited INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (personnel_id) REFERENCES personnels(personnel_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS research_authors (
    author_id SERIAL PRIMARY KEY,
    research_id INT NOT NULL,
    author_name TEXT NOT NULL,
    author_order INT NOT NULL,
    FOREIGN KEY (research_id) REFERENCES research(research_id) ON DELETE CASCADE,
    UNIQUE (research_id, author_order)
);

-- create admission

CREATE TABLE IF NOT EXISTS admission (
    admission_id SERIAL PRIMARY KEY,
    round TEXT NOT NULL,
    detail TEXT NOT NULL,
    file_image TEXT NOT NULL
);

-- create calendar

CREATE TABLE IF NOT EXISTS calendar (
    calendar_id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,  
    detail TEXT NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL
);

-- insert calendar

INSERT INTO calendar(title,detail,start_date,end_date) VALUES
('มหกรรมการสอบ','มหกรรมการสอบ', '2025-10-27', '2025-11-08'),
('test1','test1','2025-10-29', '2025-11-05'),
('test2','test2','2025-11-09', '2025-11-09');

-- create user

CREATE TABLE users (
    user_id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NULL,
    is_active BOOLEAN DEFAULT true,
    email_verified BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_login TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

-- Index สำหรับ login
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_active ON users(is_active);

INSERT INTO users(username, email) VALUES
('kaewjamnong_s','kaewjamnong_s@su.ac.th'),
('wongtaweesap_o','wongtaweesap_o@su.ac.th'),
('worrawichaipat_p','worrawichaipat_p@su.ac.th'),
('tantatsanawong_p','tantatsanawong_p@su.ac.th'),
('sitdhisanguan_k','sitdhisanguan_k@su.ac.th'),
('seepanomwan_k','seepanomwan_k@su.ac.th'),
-- ('praditwong_k','praditwong_k@su.ac.th'),
-- ('promrit_n','promrit_n@su.ac.th'),
('soonklang_t','soonklang_t@su.ac.th'),
('aonpong_p','aonpong_p@su.ac.th'),
('kanawong_r','kanawong_r@su.ac.th'),
('muangon_w','muangon_w@su.ac.th'),
-- ('waijanya_s','waijanya_s@su.ac.th'),
('pongpinigpinyo_s','pongpinigpinyo_s@su.ac.th'),
('chaowalit_o','chaowalit_o@su.ac.th'),
('tangjui_n','tangjui_n@su.ac.th'),
('pansri_b','pansri_b@su.ac.th'),
('wasara','wasara@cp.su.ac.th'), --*
('arampongsanuwat_s','arampongsanuwat_s@su.ac.th'),
('rodhetbhai_s','rodhetbhai_s@su.ac.th'),
-- ('hongwitayakorn_a','hongwitayakorn_a@su.ac.th'),
-- Admin, Staff
-- ('luangsamankul_p','luangsamankul_p@su.ac.th'),
-- ('sonsanguan_w','sonsanguan_w@su.ac.th'),
('tatong_k','tatong_k@su.ac.th'),
('jancharoen_k','jancharoen_k@su.ac.th'),
('chaysirikhae_t','chaysirikhae_t@su.ac.th'),
('yahom_s','yahom_s@su.ac.th');


-- test Admin
INSERT INTO users(username, email, password_hash) VALUES
('sonsanguan_w','sonsanguan_w@su.ac.th','$2a$12$4M9WxFsEO32LtVqOVEU37OYJ1/0Hp4cq8.X.E6KI5qP9eYNHos2ue'),
('luangsamankul_p','luangsamankul_p@su.ac.th','$2a$12$4M9WxFsEO32LtVqOVEU37OYJ1/0Hp4cq8.X.E6KI5qP9eYNHos2ue'),
('waijanya_s','waijanya_s@su.ac.th','$2a$12$4M9WxFsEO32LtVqOVEU37OYJ1/0Hp4cq8.X.E6KI5qP9eYNHos2ue'),
('promrit_n','promrit_n@su.ac.th','$2a$12$4M9WxFsEO32LtVqOVEU37OYJ1/0Hp4cq8.X.E6KI5qP9eYNHos2ue'),
('praditwong_k','praditwong_k@su.ac.th','$2a$12$4M9WxFsEO32LtVqOVEU37OYJ1/0Hp4cq8.X.E6KI5qP9eYNHos2ue'),
('hongwitayakorn_a','hongwitayakorn_a@su.ac.th','$2a$12$4M9WxFsEO32LtVqOVEU37OYJ1/0Hp4cq8.X.E6KI5qP9eYNHos2ue'),
('supapanchai_r','supapanchai_r@su.ac.th','$2a$12$4M9WxFsEO32LtVqOVEU37OYJ1/0Hp4cq8.X.E6KI5qP9eYNHos2ue'),
('saengson_s','saengson_s@su.ac.th','$2a$12$4M9WxFsEO32LtVqOVEU37OYJ1/0Hp4cq8.X.E6KI5qP9eYNHos2ue'),
('admin','admin@gmail.com','$2a$12$4M9WxFsEO32LtVqOVEU37OYJ1/0Hp4cq8.X.E6KI5qP9eYNHos2ue'),
('staff','staff@gmail.com','$2a$12$4M9WxFsEO32LtVqOVEU37OYJ1/0Hp4cq8.X.E6KI5qP9eYNHos2ue');

-- create role

CREATE TABLE roles (
    role_id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_roles_name ON roles(name);

-- insert roles

INSERT INTO roles (name, description) VALUES
('admin', 'Administrator with full system access'),
('staff', 'Can create and edit content'),
('teacher', 'Edit personal information');


CREATE TABLE user_roles (
    user_id INTEGER NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    role_id INTEGER NOT NULL REFERENCES roles(role_id) ON DELETE CASCADE,
    assigned_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    assigned_by INTEGER REFERENCES users(user_id),
    PRIMARY KEY (user_id, role_id)
);

CREATE INDEX idx_user_roles_user ON user_roles(user_id);
CREATE INDEX idx_user_roles_role ON user_roles(role_id);

INSERT INTO user_roles (user_id, role_id, assigned_by)
SELECT u.user_id, r.role_id, assigner.user_id
FROM users u
JOIN roles r ON r.name = 'teacher'
JOIN users assigner ON assigner.username = 'sonsanguan_w'
WHERE u.username IN (
    'kaewjamnong_s',
    'wongtaweesap_o',
    'worrawichaipat_p',
    'tantatsanawong_p',
    'sitdhisanguan_k',
    'seepanomwan_k',
    'praditwong_k',
    'promrit_n',
    'soonklang_t',
    'aonpong_p',
    'kanawong_r',
    'muangon_w',
    'waijanya_s',
    'pongpinigpinyo_s',
    'chaowalit_o',
    'tangjui_n',
    'pansri_b',
    'wasara',
    'arampongsanuwat_s',
    'rodhetbhai_s',
    'hongwitayakorn_a'
);

INSERT INTO user_roles (user_id, role_id, assigned_by)
SELECT u.user_id, r.role_id, assigner.user_id
FROM users u
JOIN roles r ON r.name = 'staff'
JOIN users assigner ON assigner.username = 'sonsanguan_w'
WHERE u.username IN (
    'luangsamankul_p',
    'tatong_k',
    'jancharoen_k',
    'chaysirikhae_t',
    'yahom_s',
    'staff'
);

INSERT INTO user_roles (user_id, role_id, assigned_by)
SELECT u.user_id, r.role_id, assigner.user_id
FROM users u
JOIN roles r ON r.name = 'admin'
JOIN users assigner ON assigner.username = 'sonsanguan_w'
WHERE u.username IN (
    'sonsanguan_w',
    'supapanchai_r',
    'saengson_s',
    'admin'
);

-- create permissions

CREATE TABLE permissions (
    permission_id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    resource VARCHAR(50) NOT NULL,
    action VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_permissions_name ON permissions(name);
CREATE INDEX idx_permissions_resource ON permissions(resource);

-- insert permissions

INSERT INTO permissions (name, description, resource, action) VALUES
-- news
('news:read', 'Can view news', 'news', 'read'),
('news:read_id', 'Can view news id', 'news', 'read'),
('news:create', 'Can create new news', 'news', 'create'),
('news:update', 'Can update news', 'news', 'update'),
('news:delete', 'Can delete news', 'news', 'delete'),

-- courses 
('courses:read', 'Can view courses', 'courses', 'read'),
('courses:read_id', 'Can view courses id', 'courses', 'read'),
('courses:create', 'Can create new courses', 'courses', 'create'),
('courses:update', 'Can update courses', 'courses', 'update'),
('courses:delete', 'Can delete courses', 'courses', 'delete'),

-- course_structure 
('course_structure:read', 'Can view course_structure', 'course_structure', 'read'),
('course_structure:read_id', 'Can view course_structure id', 'course_structure', 'read'),
('course_structure:create', 'Can create new course_structure', 'course_structure', 'create'),
('course_structure:update', 'Can update new course_structure', 'course_structure', 'update'),
('course_structure:delete', 'Can delete course_structure', 'course_structure', 'delete'),

-- roadmap 
('roadmap:read', 'Can view roadmap', 'roadmap', 'read'),
('roadmap:read_id', 'Can view roadmap id', 'roadmap', 'read'),
('roadmap:create', 'Can create new roadmap', 'roadmap', 'create'),
('roadmap:delete', 'Can delete roadmap', 'roadmap', 'delete'),

-- subject 
('subject:read', 'Can view subject', 'subject', 'read'),
('subject:read_id', 'Can view subject id', 'subject', 'read'),
('subject:create', 'Can create new subject', 'subject', 'create'),
('subject:update', 'Can update subject', 'subject', 'update'),
('subject:delete', 'Can delete subject', 'subject', 'delete'),

-- personnel
('personnel:read', 'Can view personnel', 'personnel', 'read'),
('personnel:read_id', 'Can view personnel id', 'personnel', 'read'),
('personnel:create', 'Can create new personnel', 'personnel', 'create'),
('personnel:update', 'Can update personnel', 'personnel', 'update'),
('your_personnel:update', 'Can update your personal information', 'personnel', 'update'),
('personnel:delete', 'Can delete personnel', 'personnel', 'delete'),

-- research
('scopus:read', 'Research data is accessible', 'research', 'read'),
('research:read', 'Can view research', 'research', 'read'),

-- admission
('admission:read', 'Can view admission', 'admission', 'read'),
('admission:read_id', 'Can view admission id', 'admission', 'read'),
('admission:create', 'Can create new admission', 'admission', 'create'),
('admission:update', 'Can update admission', 'admission', 'update'),
('admission:delete', 'Can delete admission', 'admission', 'delete'),

-- calendar
('calendar:read', 'Can view calendar', 'calendar', 'read'),
('calendar:read_id', 'Can view calendar id', 'calendar', 'read'),
('calendar:create', 'Can create new calendar', 'calendar', 'create'),
('calendar:update', 'Can update calendar', 'calendar', 'update'),
('calendar:delete', 'Can delete calendar', 'calendar', 'delete'),

-- users
('users:read', 'Can view users', 'users', 'read'),
('users:create', 'Can create new users', 'users', 'create'),
('users:delete', 'Can delete users', 'users', 'delete'),

-- Roles
('roles:assign', 'Can assign roles to users', 'roles', 'assign'),

-- logs
('logs:read', 'Can view logs', 'logs', 'read');

CREATE TABLE role_permissions (
    role_id INTEGER NOT NULL REFERENCES roles(role_id) ON DELETE CASCADE,
    permission_id INTEGER NOT NULL REFERENCES permissions(permission_id) ON DELETE CASCADE,
    granted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (role_id, permission_id)
);

CREATE INDEX idx_role_perms_role ON role_permissions(role_id);
CREATE INDEX idx_role_perms_perm ON role_permissions(permission_id);

-- insert role permissions

-- Admin: ทุก permissions
INSERT INTO role_permissions (role_id, permission_id)
SELECT
    (SELECT role_id FROM roles WHERE name = 'admin'),
    permission_id FROM permissions;

-- Staff: ทุก permissions (ยกเว้น Roles permissions)
INSERT INTO role_permissions (role_id, permission_id)
SELECT
    (SELECT role_id FROM roles WHERE name = 'staff'),
    permission_id FROM permissions
WHERE name IN (
    'news:read', 'news:read_id', 'news:create', 'news:update', 'news:delete',
    'courses:read', 'courses:read_id', 'courses:create', 'courses:update', 'courses:delete',
    'course_structure:read', 'course_structure:read_id', 'course_structure:create', 'course_structure:update', 'course_structure:delete',
    'roadmap:read', 'roadmap:read_id', 'roadmap:create', 'roadmap:delete',
    'subject:read', 'subject:read_id', 'subject:create', 'subject:update', 'subject:delete',
    'personnel:read', 'personnel:read_id', 'personnel:create', 'personnel:update', 'your_personnel:update', 'personnel:delete',
    'scopus:read', 'research:read',
    'admission:read', 'admission:read_id', 'admission:create', 'admission:update', 'admission:delete',
    'calendar:read', 'calendar:read_id', 'calendar:create', 'calendar:update', 'calendar:delete'
);

-- Teacher your_personnel:update
INSERT INTO role_permissions (role_id, permission_id)
SELECT
    (SELECT role_id FROM roles WHERE name = 'teacher'),
    permission_id FROM permissions WHERE name = 'your_personnel:update';

-- Refresh Tokens (สำหรับ JWT refresh)
CREATE TABLE refresh_tokens (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    token VARCHAR(500) UNIQUE NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    revoked_at TIMESTAMP,
    replaced_by VARCHAR(500)
);

CREATE INDEX idx_refresh_tokens_user ON refresh_tokens(user_id);
CREATE INDEX idx_refresh_tokens_token ON refresh_tokens(token);
CREATE INDEX idx_refresh_tokens_expires ON refresh_tokens(expires_at);

-- Audit Logs (สำหรับ tracking)
CREATE TABLE audit_logs (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(user_id),
    action VARCHAR(100) NOT NULL,  
    resource VARCHAR(50),           
    resource_id VARCHAR(50),
    details JSONB,
    ip_address VARCHAR(50),
    user_agent TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_audit_logs_user ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_action ON audit_logs(action);
CREATE INDEX idx_audit_logs_created ON audit_logs(created_at);

-- TRIGGER

CREATE OR REPLACE FUNCTION update_modified_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP; 
    RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';

CREATE TRIGGER update_news_modtime
BEFORE UPDATE ON news
FOR EACH ROW
EXECUTE FUNCTION update_modified_column();

CREATE TRIGGER update_users_modtime
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION update_modified_column();

CREATE TRIGGER update_roles_modtime
BEFORE UPDATE ON roles
FOR EACH ROW
EXECUTE FUNCTION update_modified_column();