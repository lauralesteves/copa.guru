import { useEffect } from 'react';
import { Link } from 'react-router-dom';
import { useLocale } from '../../i18n/LocaleContext';

export function PrivacyPolicy() {
  const { t } = useLocale();
  useEffect(() => {
    window.scrollTo(0, 0);
    document.title = 'Política de Privacidade — Copa.Guru';
  }, []);

  return (
    <div className="min-h-dvh bg-copa-dark pt-24 pb-16 px-6">
      <div className="max-w-3xl mx-auto">
        <Link
          to="/"
          className="text-copa-gold hover:text-copa-gold-light text-sm transition-colors mb-8 inline-block"
        >
          {t.legal.back}
        </Link>

        <h1 className="font-display text-4xl sm:text-5xl text-copa-gold tracking-wider mb-8">
          POLÍTICA DE PRIVACIDADE
        </h1>

        <div className="prose prose-invert max-w-none text-white/80 space-y-6 text-sm leading-relaxed">
          <h2 className="text-white text-lg font-semibold mt-8">Política de Privacidade</h2>

          <p>
            A sua privacidade é importante para nós. É política do Copa.Guru respeitar a sua
            privacidade em relação a qualquer informação sua que possamos coletar no site{' '}
            <a href="https://copa.guru" className="text-copa-gold hover:text-copa-gold-light transition-colors">Copa.Guru</a>,
            e outros sites que possuímos e operamos.
          </p>

          <p>
            Solicitamos informações pessoais apenas quando realmente precisamos delas para lhe
            fornecer um serviço. Fazemo-lo por meios justos e legais, com o seu conhecimento e
            consentimento. Também informamos por que estamos coletando e como será usado.
          </p>

          <p>
            Apenas retemos as informações coletadas pelo tempo necessário para fornecer o serviço
            solicitado. Quando armazenamos dados, protegemos dentro de meios comercialmente
            aceitáveis para evitar perdas e roubos, bem como acesso, divulgação, cópia, uso ou
            modificação não autorizados.
          </p>

          <p>
            Não compartilhamos informações de identificação pessoal publicamente ou com terceiros,
            exceto quando exigido por lei.
          </p>

          <p>
            O nosso site pode ter links para sites externos que não são operados por nós. Esteja
            ciente de que não temos controle sobre o conteúdo e práticas desses sites e não podemos
            aceitar responsabilidade por suas respectivas políticas de privacidade.
          </p>

          <p>
            Você é livre para recusar a nossa solicitação de informações pessoais, entendendo que
            talvez não possamos fornecer alguns dos serviços desejados.
          </p>

          <p>
            O uso continuado de nosso site será considerado como aceitação de nossas práticas em
            torno de privacidade e informações pessoais. Se você tiver alguma dúvida sobre como
            lidamos com dados do usuário e informações pessoais, entre em contato conosco.
          </p>

          <h2 className="text-white text-lg font-semibold mt-8">Política de Cookies</h2>

          <h3 className="text-white/90 font-semibold">O que são cookies?</h3>

          <p>
            Como é prática comum em quase todos os sites profissionais, este site usa cookies, que
            são pequenos arquivos baixados no seu computador, para melhorar sua experiência. Esta
            página descreve quais informações eles coletam, como as usamos e por que às vezes
            precisamos armazenar esses cookies.
          </p>

          <h3 className="text-white/90 font-semibold">Como usamos os cookies?</h3>

          <p>
            Utilizamos cookies por vários motivos, detalhados abaixo. Infelizmente, na maioria dos
            casos, não existem opções padrão do setor para desativar os cookies sem desativar
            completamente a funcionalidade e os recursos que eles adicionam a este site. É
            recomendável que você deixe todos os cookies se não tiver certeza se precisa ou não
            deles, caso sejam usados para fornecer um serviço que você usa.
          </p>

          <h3 className="text-white/90 font-semibold">Desativar cookies</h3>

          <p>
            Você pode impedir a configuração de cookies ajustando as configurações do seu navegador.
            Esteja ciente de que a desativação de cookies afetará a funcionalidade deste e de muitos
            outros sites que você visita. Portanto, é recomendável que você não desative os cookies.
          </p>

          <h3 className="text-white/90 font-semibold">Cookies que definimos</h3>

          <ul className="list-disc pl-6 space-y-2">
            <li>
              <strong>Cookies relacionados a preferências do site</strong> — Para proporcionar uma
              ótima experiência neste site, fornecemos a funcionalidade para definir suas
              preferências (como idioma e palpites salvos). Para lembrar suas preferências,
              precisamos definir cookies ou usar localStorage.
            </li>
          </ul>

          <h3 className="text-white/90 font-semibold">Dados armazenados localmente</h3>

          <p>
            O Copa.Guru utiliza o localStorage do navegador para salvar seus palpites e preferências.
            Esses dados permanecem apenas no seu dispositivo e não são enviados para nenhum servidor.
            Você pode limpar esses dados a qualquer momento nas configurações do seu navegador ou
            usando o botão "Limpar tudo" na seção de palpites.
          </p>

          <h3 className="text-white/90 font-semibold">Mais informações</h3>

          <p>
            Esperamos que esteja esclarecido e, como mencionado anteriormente, se houver algo que
            você não tem certeza se precisa ou não, geralmente é mais seguro deixar os cookies
            ativados, caso interaja com um dos recursos que você usa em nosso site.
          </p>

          <p className="text-white/50 mt-8">
            Esta política é efetiva a partir de <strong>Março/2026</strong>.
          </p>
        </div>
      </div>
    </div>
  );
}
