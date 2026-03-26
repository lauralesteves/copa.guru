import { useEffect } from 'react';
import { Link } from 'react-router-dom';

export function Terms() {
  useEffect(() => {
    window.scrollTo(0, 0);
    document.title = 'Termos & Condições — Copa.Guru';
  }, []);

  return (
    <div className="min-h-dvh bg-copa-dark pt-24 pb-16 px-6">
      <div className="max-w-3xl mx-auto">
        <Link
          to="/"
          className="text-copa-gold hover:text-copa-gold-light text-sm transition-colors mb-8 inline-block"
        >
          ← Voltar para Copa.Guru
        </Link>

        <h1 className="font-display text-4xl sm:text-5xl text-copa-gold tracking-wider mb-8">
          TERMOS & CONDIÇÕES
        </h1>

        <div className="prose prose-invert max-w-none text-white/80 space-y-6 text-sm leading-relaxed">
          <h2 className="text-white text-lg font-semibold mt-8">1. Termos</h2>

          <p>
            Ao acessar o site{' '}
            <a href="https://copa.guru" className="text-copa-gold hover:text-copa-gold-light transition-colors">Copa.Guru</a>,
            você concorda em cumprir estes termos de serviço, todas as leis e regulamentos
            aplicáveis, e concorda que é responsável pelo cumprimento de todas as leis locais
            aplicáveis. Se você não concordar com algum desses termos, está proibido de usar ou
            acessar este site. Os materiais contidos neste site são protegidos pelas leis de direitos
            autorais e marcas comerciais aplicáveis.
          </p>

          <h2 className="text-white text-lg font-semibold mt-8">2. Descrição do Serviço</h2>

          <p>
            O Copa.Guru é uma plataforma de visualização de dados e entretenimento relacionada à
            Copa do Mundo FIFA 2026. O site permite visualizar grupos, resultados, classificações,
            fazer palpites e comparar seleções. Os palpites são armazenados localmente no navegador
            do usuário e não constituem apostas ou jogos de azar.
          </p>

          <h2 className="text-white text-lg font-semibold mt-8">3. Uso de Licença</h2>

          <p>
            É concedida permissão para baixar temporariamente uma cópia dos materiais (informações
            ou software) no site Copa.Guru, apenas para visualização transitória pessoal e não
            comercial. Esta é a concessão de uma licença, não uma transferência de título e, sob esta
            licença, você não pode:
          </p>

          <ol className="list-decimal pl-6 space-y-1">
            <li>Modificar ou copiar os materiais;</li>
            <li>Usar os materiais para qualquer finalidade comercial ou para exibição pública (comercial ou não comercial);</li>
            <li>Tentar descompilar ou fazer engenharia reversa de qualquer software contido no site Copa.Guru;</li>
            <li>Remover quaisquer direitos autorais ou outras notações de propriedade dos materiais;</li>
            <li>Transferir os materiais para outra pessoa ou "espelhar" os materiais em qualquer outro servidor.</li>
          </ol>

          <p>
            Esta licença será automaticamente rescindida se você violar alguma dessas restrições e
            poderá ser rescindida por Copa.Guru a qualquer momento.
          </p>

          <h2 className="text-white text-lg font-semibold mt-8">4. Isenção de Responsabilidade</h2>

          <p>
            Os materiais no site do Copa.Guru são fornecidos "como estão". Copa.Guru não oferece
            garantias, expressas ou implícitas, e, por este meio, isenta e nega todas as outras
            garantias, incluindo, sem limitação, garantias implícitas ou condições de
            comercialização, adequação a um fim específico ou não violação de propriedade intelectual
            ou outra violação de direitos.
          </p>

          <p>
            Além disso, o Copa.Guru não garante ou faz qualquer representação relativa à precisão,
            aos resultados prováveis ou à confiabilidade do uso dos materiais em seu site ou de outra
            forma relacionado a esses materiais ou em sites vinculados a este site.
          </p>

          <h2 className="text-white text-lg font-semibold mt-8">5. Propriedade Intelectual</h2>

          <p>
            Os dados de partidas, seleções e competições são de propriedade da FIFA. As bandeiras e
            escudos são de domínio público ou de propriedade de suas respectivas federações. O
            Copa.Guru utiliza esses recursos para fins informativos e de entretenimento.
          </p>

          <h2 className="text-white text-lg font-semibold mt-8">6. Limitações</h2>

          <p>
            Em nenhum caso o Copa.Guru ou seus fornecedores serão responsáveis por quaisquer danos
            (incluindo, sem limitação, danos por perda de dados ou lucro ou devido a interrupção dos
            negócios) decorrentes do uso ou da incapacidade de usar os materiais em Copa.Guru.
          </p>

          <h2 className="text-white text-lg font-semibold mt-8">7. Links</h2>

          <p>
            O Copa.Guru não analisou todos os sites vinculados ao seu site e não é responsável pelo
            conteúdo de nenhum site vinculado. A inclusão de qualquer link não implica endosso por
            Copa.Guru do site. O uso de qualquer site vinculado é por conta e risco do usuário.
          </p>

          <h2 className="text-white text-lg font-semibold mt-8">8. Modificações</h2>

          <p>
            O Copa.Guru pode revisar estes termos de serviço do site a qualquer momento, sem aviso
            prévio. Ao usar este site, você concorda em ficar vinculado à versão atual desses termos
            de serviço.
          </p>

          <h2 className="text-white text-lg font-semibold mt-8">9. Lei Aplicável</h2>

          <p>
            Estes termos e condições são regidos e interpretados de acordo com as leis do Brasil e
            você se submete irrevogavelmente à jurisdição exclusiva dos tribunais naquele estado ou
            localidade.
          </p>

          <p className="text-white/50 mt-8">
            Estes termos são efetivos a partir de <strong>Março/2026</strong>.
          </p>
        </div>
      </div>
    </div>
  );
}
