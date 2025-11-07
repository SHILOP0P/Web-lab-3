--
-- PostgreSQL database dump
--

\restrict jmrBpWWIhWTjFghDcrq1vzStc861ekKKpuwG6XMLaolVea7Ub2h6beiClPWWmF5

-- Dumped from database version 18.0
-- Dumped by pg_dump version 18.0

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: product; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.product (
    id integer NOT NULL,
    manufacturer_id integer NOT NULL,
    name character varying(255) NOT NULL,
    alias character varying(255) NOT NULL,
    short_description text NOT NULL,
    description text NOT NULL,
    price numeric(20,2) NOT NULL,
    available integer DEFAULT 1 NOT NULL,
    meta_keywords character varying(255) NOT NULL,
    meta_description character varying(255) NOT NULL,
    meta_title character varying(255) NOT NULL
);


ALTER TABLE public.product OWNER TO postgres;

--
-- Name: product_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.product_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.product_id_seq OWNER TO postgres;

--
-- Name: product_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.product_id_seq OWNED BY public.product.id;


--
-- Name: product_images; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.product_images (
    id integer NOT NULL,
    product_id integer NOT NULL,
    image character varying(255) NOT NULL,
    title character varying(255) NOT NULL
);


ALTER TABLE public.product_images OWNER TO postgres;

--
-- Name: product_images_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.product_images_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.product_images_id_seq OWNER TO postgres;

--
-- Name: product_images_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.product_images_id_seq OWNED BY public.product_images.id;


--
-- Name: product_properties; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.product_properties (
    id integer NOT NULL,
    product_id integer NOT NULL,
    characteristics text
);


ALTER TABLE public.product_properties OWNER TO postgres;

--
-- Name: product_properties_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.product_properties_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.product_properties_id_seq OWNER TO postgres;

--
-- Name: product_properties_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.product_properties_id_seq OWNED BY public.product_properties.id;


--
-- Name: reviews; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.reviews (
    id integer NOT NULL,
    user_id integer NOT NULL,
    content text NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.reviews OWNER TO postgres;

--
-- Name: reviews_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.reviews_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.reviews_id_seq OWNER TO postgres;

--
-- Name: reviews_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.reviews_id_seq OWNED BY public.reviews.id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users (
    id bigint NOT NULL,
    username character varying(50) NOT NULL,
    email character varying(255) NOT NULL,
    password_hash text NOT NULL,
    role character varying(20) DEFAULT 'USER'::character varying NOT NULL,
    first_name character varying(100),
    last_name character varying(100),
    phone character varying(50),
    gender character varying(20),
    birthdate date,
    region character varying(100),
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.users OWNER TO postgres;

--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.users_id_seq OWNER TO postgres;

--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: product id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.product ALTER COLUMN id SET DEFAULT nextval('public.product_id_seq'::regclass);


--
-- Name: product_images id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.product_images ALTER COLUMN id SET DEFAULT nextval('public.product_images_id_seq'::regclass);


--
-- Name: product_properties id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.product_properties ALTER COLUMN id SET DEFAULT nextval('public.product_properties_id_seq'::regclass);


--
-- Name: reviews id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.reviews ALTER COLUMN id SET DEFAULT nextval('public.reviews_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Data for Name: product; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.product (id, manufacturer_id, name, alias, short_description, description, price, available, meta_keywords, meta_description, meta_title) FROM stdin;
13	1	Цемент М500	Строительные материалы	Цемент М500 — универсальный высокопрочный строительный материал, применяемый для изготовления бетонов, железобетонных изделий, строительных растворов и штукатурок. Отличается повышенной устойчивостью к влаге и морозу, а также долговечностью и прочностью на сжатие.	Цемент М500 используется для заливки фундаментов, плит перекрытий, дорожных покрытий, лестничных маршей и других конструкций, требующих высокой прочности и устойчивости к нагрузкам. Он также применяется при строительстве гидротехнических сооружений, благодаря отличной влагостойкости.\n\nВ состав цемента входят клинкер, гипс и минеральные добавки, что обеспечивает оптимальную скорость твердения и долговечность материала. При правильном хранении цемент сохраняет свои свойства в течение длительного времени. Материал подходит для эксплуатации в условиях сурового климата.\n\nРекомендуется использовать в сочетании с песком и щебнем в соотношениях, указанных производителем. Перед применением убедитесь, что поверхность очищена от пыли, масел и загрязнений для обеспечения максимальной адгезии.	1100.00	26	cement		cement
14	1	Гипс строительный	Строительные материалы	Гипс строительный — это экологически чистый материал, широко применяемый для внутренних отделочных работ, а также для создания декоративных элементов и лепнины. Он обеспечивает ровную поверхность и идеально подходит для оштукатуривания стен и потолков в помещениях с нормальной влажностью.	Материал отличается высокой чистотой состава, мелкодисперсной структурой и оптимальной пластичностью, что обеспечивает лёгкость нанесения и ровную поверхность. После высыхания гипсовый слой не даёт усадки и не трескается. Может применяться для создания форм, панно, декоративных элементов и потолочных плит.\n\nГипс обладает естественной паропроницаемостью, создавая здоровый микроклимат в помещении. Благодаря своей экологичности и безопасности, материал часто используется при отделке жилых домов, школ и медицинских учреждений.	950.00	18	gypsum	гипс	gypsum
21	3	Рулетка 5м	Инструменты	Надёжная измерительная рулетка длиной 5 метров с прочным корпусом и блокировкой полотна. Идеально подходит для строительных, ремонтных и монтажных работ.	Корпус рулетки покрыт резиновыми вставками, предотвращающими скольжение и защищающими инструмент от ударов при падении. Контрастная шкала легко читается при любом освещении, а фиксатор надёжно удерживает ленту в нужном положении.\n\nБлагодаря компактным размерам и малому весу рулетка удобно лежит в руке и легко помещается в карман или сумку. Это незаменимый инструмент для любого мастера.	300.00	22	tape	tape	tape
22	3	Шпатель 100мм	Инструменты	Универсальный шпатель шириной 100 мм — инструмент для нанесения, выравнивания и удаления строительных смесей. Подходит для штукатурных, шпатлёвочных и малярных работ.	Лезвие из нержавеющей стали обеспечивает долгий срок службы и устойчивость к коррозии. Благодаря гибкости металла шпатель равномерно распределяет раствор и помогает добиться идеальной гладкости поверхности.\n\nИнструмент подходит для работы с различными типами смесей: цементными, гипсовыми и полимерными. Легко очищается от остатков материала и не требует особого ухода.	180.00	29	spatula	spatula	spatula
15	1	Клей плиточный	Строительные материалы	Плиточный клей предназначен для надёжного крепления керамической, каменной и мозаичной плитки на поверхности из бетона, цемента, кирпича и гипсокартона. Обеспечивает прочное сцепление, устойчив к влаге и перепадам температур, подходит как для внутренних, так и для наружных работ.	Клей обеспечивает пластичность и лёгкость нанесения, не сползает при укладке плитки на вертикальные поверхности. После высыхания образует прочный и долговечный слой, устойчивый к влаге и механическим нагрузкам. Подходит для применения в ванных комнатах, кухнях, коридорах и на балконах.\n\nМатериал можно использовать для систем "тёплый пол". Благодаря своей структуре, клей компенсирует температурные расширения и предотвращает растрескивание плиточного покрытия. Подходит для плит среднего и крупного формата.	780.00	30	glue	glue	glue
17	2	Краска фасадная	Отделочные материалы	Фасадная краска предназначена для окрашивания внешних стен зданий, цоколей и архитектурных элементов. Обеспечивает стойкое покрытие, устойчивое к атмосферным осадкам, ультрафиолету и резким перепадам температуры. Подходит для бетона, кирпича, штукатурки и цементных поверхностей.	Фасадная краска образует прочную и эластичную плёнку, которая надёжно защищает поверхность от влаги и грязи. Она легко наносится кистью, валиком или краскопультом, не образует подтёков и быстро высыхает. Подходит для реставрационных и отделочных работ, а также для окрашивания фасадов жилых домов и общественных зданий.\n\nБлагодаря высокой укрывистости и адгезии, краска обеспечивает долговечное и эстетичное покрытие. Может использоваться как в новом строительстве, так и при ремонте старых фасадов.	979.99	18	paint	paint	paint
18	2	Обои виниловые	Отделочные материалы	Виниловые обои на флизелиновой основе — универсальный материал для внутренней отделки помещений. Отличаются прочностью, влагостойкостью и устойчивостью к ультрафиолету. Обеспечивают ровную и гладкую поверхность, скрывая мелкие дефекты стен.	Виниловые обои легко клеятся, не растягиваются и не рвутся при нанесении. Они долговечны и устойчивы к воздействию влаги, поэтому отлично подходят для кухни, прихожей и ванной комнаты. Благодаря плотной структуре, материал скрывает неровности стен и обеспечивает дополнительную теплоизоляцию.\n\nОбои можно окрашивать до 5 раз без потери текстуры. Они не выгорают на солнце, сохраняют насыщенность цвета и не требуют специального ухода. Для поклейки рекомендуется использовать клей для тяжёлых обоев.	1200.00	18	wallpaper	wallpaper	wallpaper
19	2	Плитка керамическая	Отделочные материалы	Керамическая плитка — долговечный и надёжный облицовочный материал, идеально подходящий для отделки стен и полов. Отличается устойчивостью к влаге, механическим повреждениям и воздействию химических веществ. Применяется в ванных, кухнях, санузлах и общественных помещениях.	Плитка обладает высокой прочностью и устойчивостью к износу, легко очищается и сохраняет внешний вид в течение многих лет. Может использоваться в помещениях с повышенной влажностью и перепадами температур. Благодаря современному дизайну и разнообразию оттенков, плитка создаёт стильный и аккуратный интерьер.\n\nДля укладки рекомендуется использовать плиточный клей с высокой адгезией и влагостойкий затирочный состав. Поверхность перед монтажом должна быть ровной и сухой. При правильном уходе плитка прослужит десятилетиями.	1350.00	27	tile	tile	tile
20	3	Дрель Bosch	Инструменты	Электрическая дрель Bosch — надёжный инструмент для сверления металла, дерева, пластика и других материалов. Отличается компактностью, мощным двигателем и эргономичным дизайном, обеспечивающим удобство при работе.	Дрель Bosch подходит как для бытового, так и профессионального применения. Эргономичный корпус с прорезиненной рукояткой обеспечивает надёжный захват, а встроенная вентиляция защищает двигатель от перегрева. Быстрозажимной патрон позволяет легко менять сверла без дополнительных инструментов.\n\nИнструмент может использоваться для сверления в режиме шуруповёрта, имеет точный регулятор оборотов и устойчиво работает при длительных нагрузках. Отличное сочетание цены, мощности и качества.	4500.00	17	drill	drill	drill
\.


--
-- Data for Name: product_images; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.product_images (id, product_id, image, title) FROM stdin;
16	13	1761840556441779400_.jpg	cement.jpg
17	14	1761847209504133700_.jpg	gypsum.jpg
18	15	1761847337731219300_.jpg	glue.jpg
19	17	1761847699325857900_.jpg	paint.jpg
20	18	1761847906146647600_.jpg	wallpaper.jpg
21	19	1761848008205872600_.jpg	tile.jpg
22	20	1761848116986901500_.jpg	drill.jpg
23	21	1761848222190823400_.jpg	tape.jpg
24	22	1761848310776399200_.jpg	spatula.jpg
\.


--
-- Data for Name: product_properties; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.product_properties (id, product_id, characteristics) FROM stdin;
4	13	Марка прочности: М500 Д0\r\nФасовка: мешок 50 кг\r\nЦвет: серый\r\nСрок схватывания: не менее 45 минут\r\nПрочность на сжатие: ≥ 50 МПа\r\nМорозостойкость: F300\r\nВодостойкость: W6\r\nСтрана-производитель: Россия
5	14	Тип: Гипс строительный Г-5 Б II\r\nФасовка: мешок 30 кг\r\nЦвет: белый / светло-серый\r\nВремя схватывания: 6–15 мин\r\nПрочность на сжатие: не менее 5 МПа\r\nТемпература применения: от +5 до +30 °C\r\nСтрана-производитель: Россия
6	15	Тип: цементный клей С0/С1\r\nФасовка: мешок 25 кг\r\nРасход: 3–5 кг/м²\r\nВремя корректировки плитки: до 15 минут\r\nТемпература нанесения: от +5 до +35 °C\r\nАдгезия: ≥ 0.5 МПа\r\nМорозостойкость: не менее 50 циклов
8	17	Тип: акриловая водоэмульсионная краска\r\nФасовка: ведро 10 л\r\nРасход: 0.18–0.25 л/м²\r\nВремя высыхания: 1–2 часа\r\nЦвет: белый (возможна колеровка)\r\nСтойкость к УФ-излучению: высокая\r\nСрок службы покрытия: до 10 лет
9	18	Тип: виниловые на флизелиновой основе\r\nРазмер рулона: 1.06 × 10 м\r\nПлотность: 220 г/м²\r\nУстойчивость к влаге: высокая\r\nСветостойкость: отличная\r\nВозможность окрашивания: да\r\nСтрана-производитель: Германия
10	19	Тип: настенная / напольная плитка\r\nРазмер: 300×300 мм\r\nТолщина: 8 мм\r\nПокрытие: глянцевое\r\nИзносостойкость: класс PEI III\r\nВодопоглощение: ≤ 0.5%\r\nСтрана-производитель: Испания
11	20	Мощность: 600 Вт\r\nТип патрона: быстрозажимной, 13 мм\r\nЧастота вращения: 0–3000 об/мин\r\nФункция реверса: есть\r\nВес: 1.7 кг\r\nДлина кабеля: 2 м\r\nПроизводитель: Германия
12	21	Длина ленты: 5 м\r\nШирина ленты: 19 мм\r\nМатериал корпуса: ударопрочный пластик\r\nМеханизм фиксации: кнопка блокировки\r\nПоясное крепление: есть\r\nПокрытие шкалы: износостойкий лак\r\nПроизводитель: Китай
13	22	Ширина лезвия: 100 мм\r\nМатериал лезвия: нержавеющая сталь\r\nРукоятка: пластиковая, эргономичная\r\nТолщина лезвия: 0.6 мм\r\nВес: 80 г\r\nПроизводитель: Россия
\.


--
-- Data for Name: reviews; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.reviews (id, user_id, content, created_at) FROM stdin;
4	2	ку	2025-11-06 19:36:47.030852
10	3	привет	2025-11-07 13:28:35.153423
11	4	привет	2025-11-07 13:28:51.902711
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.users (id, username, email, password_hash, role, first_name, last_name, phone, gender, birthdate, region, created_at, updated_at) FROM stdin;
2	SHILOP0P	muxa4ev.dim@yandex.ru	$2a$12$GjuUKf0OIz0oDtO3BC/G2uIIGYSwbu1t76DIbZoWX20B60FkrigS2	USER	Dmitry	Muckhachev	89687865001	Мужской	2025-11-06	Москва и Московская область	2025-11-06 15:29:13.39706+03	2025-11-06 15:29:13.39706+03
3	User1	sdczxc@yandex.ru	$2a$12$.fyeOu8ZXO6QiSzEa4ttCOzzZSCM2WWdY5AZrKO7A2/v3/Z6ePuku	USER	Пользователь1	Пользователь1	81234567890	Мужской	2025-11-06	Санкт-Петербург	2025-11-06 19:40:03.307954+03	2025-11-06 19:40:03.307954+03
4	User2	sdkn@yandex.ru	$2a$12$AIhQaCrHAbDBF0OB9lH1SeRbl30ZZ4Ijy8CocdZ83LY.06V8CczzS	USER	Пользователь2	Пользователь2	80987654321	Женский	2025-11-01	Сибирский федеральный округ	2025-11-06 19:42:44.430434+03	2025-11-06 19:42:44.430434+03
\.


--
-- Name: product_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.product_id_seq', 22, true);


--
-- Name: product_images_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.product_images_id_seq', 24, true);


--
-- Name: product_properties_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.product_properties_id_seq', 13, true);


--
-- Name: reviews_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.reviews_id_seq', 12, true);


--
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.users_id_seq', 4, true);


--
-- Name: product_images product_images_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.product_images
    ADD CONSTRAINT product_images_pkey PRIMARY KEY (id);


--
-- Name: product product_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.product
    ADD CONSTRAINT product_pkey PRIMARY KEY (id);


--
-- Name: product_properties product_properties_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.product_properties
    ADD CONSTRAINT product_properties_pkey PRIMARY KEY (id);


--
-- Name: reviews reviews_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.reviews
    ADD CONSTRAINT reviews_pkey PRIMARY KEY (id);


--
-- Name: users users_email_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: users users_username_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_username_key UNIQUE (username);


--
-- Name: idx_reviews_created_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_reviews_created_at ON public.reviews USING btree (created_at DESC);


--
-- Name: idx_users_email; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_users_email ON public.users USING btree (email);


--
-- Name: idx_users_username; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_users_username ON public.users USING btree (username);


--
-- Name: product_images fk_imgs_product; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.product_images
    ADD CONSTRAINT fk_imgs_product FOREIGN KEY (product_id) REFERENCES public.product(id) ON DELETE CASCADE;


--
-- Name: product_images fk_product_images; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.product_images
    ADD CONSTRAINT fk_product_images FOREIGN KEY (product_id) REFERENCES public.product(id) ON DELETE CASCADE;


--
-- Name: product_properties fk_product_properties; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.product_properties
    ADD CONSTRAINT fk_product_properties FOREIGN KEY (product_id) REFERENCES public.product(id) ON DELETE CASCADE;


--
-- Name: product_properties fk_props_product; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.product_properties
    ADD CONSTRAINT fk_props_product FOREIGN KEY (product_id) REFERENCES public.product(id) ON DELETE CASCADE;


--
-- Name: reviews reviews_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.reviews
    ADD CONSTRAINT reviews_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

\unrestrict jmrBpWWIhWTjFghDcrq1vzStc861ekKKpuwG6XMLaolVea7Ub2h6beiClPWWmF5

